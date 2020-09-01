package example

import (
	"fmt"
	"github.com/fbelisk/aeolus/poll"
	"github.com/fbelisk/aeolus/util"
	"net"
)

func TcpServerRun() {
	listener, err := net.Listen("tcp", ":10101")
	fmt.Println("new tcp server")
	if err != nil {
		fmt.Println("HOLY SHIT IT DOES NOT RUN SUCCESS")
	}
	w := &util.WaitGroupWrapper{}
	poller, err := poll.Create()
	if err != nil {
		fmt.Println("poller create error")
		return
	}
	fmt.Println("create poller success")
	w.Wrap(func() {
		fmt.Println("start poller wait")
		err = poller.Wait(ConnHandle)
		if err != nil {
			fmt.Println("poller wait error")
			return
		}
	})
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("accept conn err "+ err.Error())
			continue
		}
		fmt.Println("start poller wait")
		file, err := conn.(*net.TCPConn).File()
		if err != nil {
			continue
		}
		err = poller.AddRead(int(file.Fd()))
		if err != nil {
			continue
		}
		fmt.Println("add read fd", file.Fd())
	}
	w.Wait()
}

func ConnHandle(fd int, event uint32) error {
	fmt.Printf("fd %d \n", fd)
	return nil
}