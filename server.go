package aeolus

import (
	"github.com/fbelisk/aeolus/internal"
	"github.com/fbelisk/aeolus/listener"
	"github.com/fbelisk/aeolus/protocol"
	"github.com/fbelisk/aeolus/util"
	"golang.org/x/sys/unix"
	"runtime"
	"sync"
)

type MessageHandle func(message []byte) error

type Server struct {
	ln *listener.Listener
	messageHandle MessageHandle
	len       int
	lastIndex int
	Pollers   []*internal.Poller
}

func Run(p protocol.Protocal, handle MessageHandle, addr string) error {
	sync.WaitGroup{}
	ln, err := listener.NewListener(p, addr)
	if err != nil {
		return err
	}

	s := &Server{
		ln: ln,
	}

	num := runtime.NumCPU()
	w := util.WaitGroupWrapper{}
	for i := 0; i < num; i++ {
		p, err := internal.Create()
		if err != nil {
			return err
		}
		s.Pollers[i] = p
		w.Wrap(func() {
			err = p.Wait(s.PollerHandle)
			if err != nil {
				//todo log
			}
		})
	}
	err = s.NextPoller().AddRead(ln.Fd)
	if err != nil {
		return err
	}

	if err != nil {
		return err
	}
}

func (s *Server) PollerHandle(fd int, event uint32) error {
	switch fd {
	case s.ln.Fd:
		nfd, _, err := unix.Accept(fd)
		if err != nil {
			if err == unix.EAGAIN {
				return nil
			}
			return err
		}
		if err := unix.SetNonblock(nfd, true); err != nil {
			return err
		}

		//add connect fd into poller
		err = s.NextPoller().AddRead(nfd)
		if err != nil {
			return err
		}

		return nil
	default:
		s.messageHandle()
	}

}

func (s *Server) NewEventLoop(handle internal.Handler) (*EventLoop, error) {
	return s, nil
}

//todo lock multi goroutine && more load Balance
func (s *Server) NextPoller() *internal.Poller {
	nextIndex := 0
	if s.lastIndex < s.len {
		nextIndex = s.lastIndex + 1
		s.lastIndex = nextIndex
	}
	return s.Pollers[nextIndex]
}
