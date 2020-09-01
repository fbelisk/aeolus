package eventpoller

import (
	"fmt"
	"github.com/fbelisk/aeolus/poll"
	"net"
)

type Eventpoller struct {
	p       *poll.Poller
	clients map[int]*net.Conn
}

func CreatePoller() (*Eventpoller, error) {
	var err error
	ep := new(Eventpoller)
	ep.clients = make(map[int]*net.Conn)
	ep.p, err = poll.Create()
	if err != nil {
		fmt.Println("poller create error")
		_ = ep.p.Close()
		return nil, err
	}
	return ep, nil
}

func (ep *Eventpoller) Run(handler poll.Handler) error {
	fmt.Println("start poller wait")
	err := ep.p.Wait(handler)
	if err != nil {
		fmt.Println("poller wait error")
		return err
	}
}
