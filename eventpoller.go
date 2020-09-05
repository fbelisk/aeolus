package sophonn

import (
	"fmt"
	"github.com/fbelisk/aeolus/poll"
	"golang.org/x/sys/unix"
)

type Eventpoller struct {
	p          *poll.Poller
	clients    map[int]*Conn
	React      func(inframe []byte) []byte
	ReadBuffer []byte
}

func CreatePoller() (*Eventpoller, error) {
	var err error
	ep := new(Eventpoller)
	ep.clients = make(map[int]*Conn)
	ep.p, err = poll.Create()
	if err != nil {
		fmt.Println("poller create error")
		_ = ep.p.Close()
		return nil, err
	}
	return ep, nil
}

func (e *Eventpoller) Run() error {
	fmt.Println("start poller wait")
	err := e.p.Wait(func(fd int, event uint32) error {
		conn, ok := e.clients[fd]
		if !ok {
			//todo fd 关闭逻辑
			return nil
		}
		if event&poll.ReadEvents > 0 {
			//todo read loop error 处理
			return e.Read(conn)
		} else if event&poll.WriteEvents > 0 {
			return e.Write(conn)
		}
		return nil
	})
	if err != nil {
		fmt.Println("poller wait error")
		return err
	}
	return nil
}

//read handle
func (e *Eventpoller) Read(c *Conn) error {
	for {
		n, err := unix.Read(c.fd, e.ReadBuffer)
		if err != nil {
			fmt.Println("read length error " + err.Error())
			break
		}
		//解码字节存入 conn 临时 buffer
		c.buffer = c.Compressor.Decode(e.ReadBuffer[:n])
		for {
			iframe, err := c.Codec.Read(c.buffer)
			if err != nil {
				return err
			}
			if iframe == nil {
				break
			}

			//业务处理
			out := e.React(iframe)
			//响应写入
			//缓冲区非空，拼接缓冲区数据和当前响应，一同写入，保证响应数据按序到达
			if !c.outBuffer.IsEmpty() {
				_, _ = c.outBuffer.Write(out)
				continue
			} else  {
				//缓冲区为空，写入，未写入部分缓存到缓冲区
				n, err := unix.Write(c.fd, out)
				if err !=nil {
					if err == unix.EAGAIN {
						_, _ = c.outBuffer.Write(out)
						_ = e.p.ModReadAndWrite(c.fd)
						continue
					}
					//todo conn close
				}
				if len(out) == n {
					continue
				}
				_, _ = c.outBuffer.Write(out[n:])
				_ = e.p.ModReadAndWrite(c.fd)
			}
		}
	}
	return nil
}

//write handle
func (e *Eventpoller) Write(c *Conn) error {
	if c.outBuffer.IsEmpty() {
		return nil
	}
	outNew, _ := c.outBuffer.LazyReadAll()
	n, err := unix.Write(c.fd, outNew)
	if err != nil {
		if err == unix.EAGAIN {
			return nil
		}
		//todo close connn
	}
	if len(outNew) == n {
		_ = e.p.ModRead(c.fd)
	}
	c.outBuffer.Shift(n)
	return nil
}
