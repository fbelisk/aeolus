package sophonn

import (
	"github.com/fbelisk/aeolus/ringbuffer"
)

type Conn struct {
	fd         int
	inBuffer   *ringbuffer.RingBuffer	//todo ringbuffer  自己实现
	outBuffer  *ringbuffer.RingBuffer	//todo 输出缓冲区数据结构待定
	buffer     []byte
}

func NewConn(fd int) *Conn {
	return &Conn{
		fd:         fd,
		inBuffer:   ringbuffer.New(0x1000),
		outBuffer:  ringbuffer.New(0x1000),
		buffer:     make([]byte, 0x1000),
	}
}
