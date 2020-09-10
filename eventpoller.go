package sophonn

import (
	"fmt"
	"github.com/fbelisk/aeolus/poll"
	"golang.org/x/sys/unix"
	"net"
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
	ep.Compressor = &CompressorDemo{}     //todo 数据压缩
	ep.Codec = &CodecDemo{}               //todo 数据编解码
	ep.ReadBuffer = make([]byte, 0x10000) //todo readbuffer 大小确定
	ep.p, err = poll.Create()
	if err != nil {
		fmt.Println("poller create error")
		_ = ep.p.Close()
		return nil, err
	}
	return ep, nil
}

func (e *Eventpoller) Run(listener *net.TCPListener) error {
	fmt.Println("start poller wait")
	//add listener
	file, err := listener.File()
	if err != nil {
		return err
	}
	_ = e.p.AddRead(int(file.Fd()))
	//poller wait
	err = e.p.Wait(func(fd int, event uint32) error {
		fmt.Printf("new epoll wait loop fd %d; evnet %d \n\r", fd, event)
		conn, ok := e.clients[fd]
		if !ok {
			return e.Accept(fd)
		}
		if event&poll.WriteEvents > 0 {
			return e.Write(conn)
		} else if event&poll.ReadEvents > 0 {
			return e.Read(conn)
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
	n, err := unix.Read(c.fd, e.ReadBuffer)
	if err != nil {
		//todo 是否判断eagain
		_ = e.Close(c)
		fmt.Println("read length error " + err.Error())
		return err
	}
	fmt.Println("recive msg " + string(e.ReadBuffer[0:n]))
	//解码字节存入 conn 临时 buffer
	//c.buffer = e.Compressor.Decode(e.ReadBuffer[:n])
	//iframe := e.Codec.Read(c.buffer)
	//if iframe == nil {
	//	break
	//}
	//todo 分包 包含反序列化、解压缩

	//业务处理
	out := e.React(e.ReadBuffer[:n])
	//响应写入
	//缓冲区非空，拼接缓冲区数据和当前响应，一同写入，保证响应数据按序到达
	if !c.outBuffer.IsEmpty() {
		fmt.Printf("readloop write to outbuffer data %s \n\r", string(out))
		_, _ = c.outBuffer.Write(out)
		return nil
	} else {
		//缓冲区为空，写入，未写入部分缓存到缓冲区
		n, err := unix.Write(c.fd, out)
		if err != nil {
			if err == unix.EAGAIN {
				fmt.Printf("readloop write eagin error fd %d data %s \n\r", c.fd, string(out))
				_, _ = c.outBuffer.Write(out)
				_ = e.p.ModReadAndWrite(c.fd)
				return nil
			}
			_ = e.Close(c)
			return err
		}
		fmt.Printf("readloop write sucess fd %d remain_len %d remain_data %s \n\r", c.fd, len(out) - n, string(out[n:]))
		if len(out) == n {
			return nil
		}
		_, _ = c.outBuffer.Write(out[n:])
		_ = e.p.ModReadAndWrite(c.fd)
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
			fmt.Printf("write eagin error fd %d data %s \n\r", c.fd, string(outNew))
			return nil
		}
		_ = e.Close(c)
	}
	if len(outNew) == n {
		_ = e.p.ModRead(c.fd)
	}
	c.outBuffer.Shift(n)
	fmt.Printf("write success fd %d data %s \n\r", c.fd, string(outNew[:n]))
	return nil
}

//conn accept
func (e *Eventpoller) Accept(fd int) error {
	nfd, _, err := unix.Accept(fd)
	if err != nil {
		return err
	}
	fmt.Printf("accept new conn fd %d \n\r", nfd)
	if conn, ok := e.clients[nfd]; ok {
		_ = e.Close(conn)
	}
	e.clients[nfd] = NewConn(nfd)
	_ = e.p.AddRead(nfd)
	return nil
}

//conn close
func (e Eventpoller) Close(c *Conn) error {
	fmt.Printf("close conn fd:%d \n\r", c.fd)
	_ = e.p.Delete(c.fd)
	delete(e.clients, c.fd)
	_ = unix.Close(c.fd)
	return nil
}
