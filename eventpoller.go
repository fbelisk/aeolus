package sophonn

import (
	"fmt"
	"github.com/fbelisk/aeolus/poll"
	"golang.org/x/sys/unix"
)

type BusinessHandler func(inframe []byte) []byte

type Eventpoller struct {
	p          *poll.Poller
	clients    map[int]*Conn
	React      BusinessHandler
	ReadBuffer []byte
	Compressor Compressor
	Codec      Codec
}

func CreatePoller(react BusinessHandler) (*Eventpoller, error) {
	var err error
	ep := new(Eventpoller)
	ep.clients = make(map[int]*Conn)
	ep.React = react
	ep.ReadBuffer = make([]byte, 0x10000)	//todo readbuffer 大小确定
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
			return e.Read(conn)
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
		c.buffer = e.Compressor.Decode(e.ReadBuffer[:n])
		for {
			iframe, err := e.Codec.Read(c.buffer)
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
					fmt.Println("我的天 connection 异常")
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
		fmt.Println("我的天 connection 异常")
	}
	if len(outNew) == n {
		_ = e.p.ModRead(c.fd)
	}
	c.outBuffer.Shift(n)
	return nil
}

//conn accept
func (e *Eventpoller) Accept(fd int) error{
	if conn, ok := e.clients[fd]; ok {
		_ = e.Close(conn)
	}
	e.clients[fd] = NewConn(fd)
	_ = e.p.AddRead(fd)
	return nil
}

//conn close
func (e Eventpoller) Close(c *Conn) error {
	//todo
	return nil
}