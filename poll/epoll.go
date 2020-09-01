// +build linux

package poll

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
	fd    int //epoll fd
	jobfd int //job event fd
}

type Handler func(fd int, event uint32) error

func Create() (poll *Poller, err error) {
	poller := new(Poller)
	poller.fd, err = unix.EpollCreate1(unix.EPOLL_CLOEXEC)
	if err != nil {
		_ = poller.Close()
		return nil, err
	}

	if poller.jobfd, err = unix.Eventfd(0, unix.EFD_NONBLOCK|unix.EFD_CLOEXEC); err != nil {
		return nil, err
	}
	return poller, nil
}

func (p *Poller) Close() error {
	if err := unix.Close(p.fd); err != nil {
		return err
	}
	return unix.Close(p.jobfd)
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
	var wakenUp bool
	for {
		n, err := unix.EpollWait(p.fd, events.eventList, -1)
		if err != nil && err != syscall.EINTR{
			log.Println(err)
			continue
		}
		for i := 1; i <= n; i++ {
			if int(events.eventList[i].Fd) == p.jobfd {
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
