// +build linux

package poll

import (
	"golang.org/x/sys/unix"
	"log"
	"syscall"
)

const (
	ReadEvents      = unix.EPOLLPRI | unix.EPOLLIN
	WriteEvents     = unix.EPOLLOUT
	ReadWriteEvents = ReadEvents | WriteEvents
)

type Poller struct {
	fd    int //epoll fd
	Jobfd int //job event fd
}

type Handler func(fd int, event uint32) error

func Create() (poll *Poller, err error) {
	poller := new(Poller)
	poller.fd, err = unix.EpollCreate1(unix.EPOLL_CLOEXEC)
	if err != nil {
		_ = poller.Close()
		return nil, err
	}

	if poller.Jobfd, err = unix.Eventfd(0, unix.EFD_NONBLOCK|unix.EFD_CLOEXEC); err != nil {
		return nil, err
	}
	return poller, nil
}

func (p *Poller) Close() error {
	if err := unix.Close(p.fd); err != nil {
		return err
	}
	return unix.Close(p.Jobfd)
}

func (p *Poller) AddRead(fd int) (err error) {
	return unix.EpollCtl(p.fd, unix.EPOLL_CTL_ADD, fd, &unix.EpollEvent{Fd: int32(fd), Events: ReadEvents})
}

func (p *Poller) AddWrite(fd int) (err error) {
	return unix.EpollCtl(p.fd, unix.EPOLL_CTL_ADD, fd, &unix.EpollEvent{Fd: int32(fd), Events: WriteEvents})
}

func (p *Poller) AddReadAndWrite(fd int) (err error) {
	return unix.EpollCtl(p.fd, unix.EPOLL_CTL_ADD, fd, &unix.EpollEvent{Fd: int32(fd), Events: ReadWriteEvents})
}

func (p *Poller) ModReadAndWrite(fd int) (err error) {
	return unix.EpollCtl(p.fd, unix.EPOLL_CTL_MOD, fd, &unix.EpollEvent{
		Events: ReadWriteEvents,
		Fd:     int32(fd),
	})
}

func (p Poller) Wait(handler Handler) (err error) {
	events := newEvents(64)
	var wakenUp bool
	for {
		n, err := unix.EpollWait(p.fd, events.eventList, -1)
		if err != nil && err != syscall.EINTR{
			log.Println(err)
			continue
		}
		for i := 1; i <= n; i++ {
			if int(events.eventList[i].Fd) == p.Jobfd {
				wakenUp = true
				log.Println("该执行任务了")
				//_, _ = unix.Read(p.wfd, p.wfdBuf)
			} else {
				err = handler(int(events.eventList[i].Fd), events.eventList[i].Events)
				log.Println(err.Error())
			}
		}

		if wakenUp {
			wakenUp = false
			//todo job queue exec
		}
		events.resizing(n)
	}
}
