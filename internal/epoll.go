// +build linux

package internal

import (
	"golang.org/x/sys/unix"
	"log"
	"syscall"
)

const (
	readEvents      = unix.EPOLLPRI | unix.EPOLLIN
	writeEvents     = unix.EPOLLOUT
	readWriteEvents = readEvents | writeEvents
)

type Poller struct {
	fd  int //epoll fd
	jfd int //job event fd
}

type Handler func(fd int, event uint32) error

func Create() (poll *Poller, err error) {
	poller := new(Poller)
	poller.fd, err = unix.EpollCreate1(unix.EPOLL_CLOEXEC)
	if err != nil {
		return nil, err
	}

	poll.jfd, err = unix.Eventfd(0, unix.EFD_NONBLOCK|unix.EFD_CLOEXEC)
	if err != nil {
		return nil, err
	}
	return poll, nil
}

func (p *Poller) Close() error {
	if err := unix.Close(p.fd); err != nil {
		return err
	}
	return unix.Close(p.jfd)
}

func (p *Poller) AddRead(fd int) (err error) {
	return unix.EpollCtl(p.fd, unix.EPOLL_CTL_ADD, fd, &unix.EpollEvent{Fd: int32(fd), Events: readEvents})
}

func (p *Poller) AddWrite(fd int) (err error) {
	return unix.EpollCtl(p.fd, unix.EPOLL_CTL_ADD, fd, &unix.EpollEvent{Fd: int32(fd), Events: writeEvents})
}

func (p *Poller) AddReadAndWrite(fd int) (err error) {
	return unix.EpollCtl(p.fd, unix.EPOLL_CTL_ADD, fd, &unix.EpollEvent{Fd: int32(fd), Events: readWriteEvents})
}

func (p Poller) Wait(handler Handler) (err error) {
	events := newEvents(64)
	for {
		n, err := unix.EpollWait(p.fd, events.eventList, -1)
		if err != nil && err != syscall.EINTR{
			log.Println(err)
			continue
		}
		for i := 1; i <= n; i++ {
			if int(events.eventList[i].Fd) == p.jfd {

			} else {
				handler(int(events.eventList[i].Fd), events.eventList[i].Events)
			}
		}

		events.resizing(n)
	}
}
