package sophonn

import (
	"github.com/fbelisk/aeolus/ringbuffer"
)

type Conn struct {
	fd         int
	inBuffer   *ringbuffer.RingBuffer	//todo ringbuffer  自己实现
	outBuffer  *ringbuffer.RingBuffer	//todo 输出缓冲区数据结构待定
	Compressor Compressor
	Codec      Codec
	buffer     []byte
}
