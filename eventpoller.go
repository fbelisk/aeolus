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
			//todo
			out := e.React(iframe)
			//响应写入
			if !c.outBuffer.IsEmpty() {
				_, _ = c.outBuffer.Write(out)
			}
			n, err := unix.Write(c.fd, c.outBuffer.Bytes())
			//io写入未就绪,写入out buffer暂存
			if err == unix.EAGAIN {
				_ = e.p.ModReadAndWrite(c.fd)
				continue
			}
			//todo buffer 需完善
			_, _ = c.outBuffer.shift(n)
			continue
		}
	}
	return nil
}

func (e *Eventpoller) Write(conn *Conn) error {

}
