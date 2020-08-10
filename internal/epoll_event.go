//+build linux

package internal

import "golang.org/x/sys/unix"

type eventList struct {
	size   int
	eventList []unix.EpollEvent
}

func newEvents(size int) *eventList {
	return &eventList{size, make([]unix.EpollEvent, size)}
}

func (el *eventList) increase() {
	el.size <<= 1
	el.eventList = make([]unix.EpollEvent, el.size)
}

//todo more resizing strategy
func (el *eventList) resizing(n int) {
	if el.size == n {
		el.increase()
	}
}