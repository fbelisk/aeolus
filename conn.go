package sophonn

import (
	"github.com/smallnest/ringbuffer"
)

type Conn struct {
	fd         int
	inBuffer   *ringbuffer.RingBuffer
	outBuffer  *ringbuffer.RingBuffer	//todo 输出缓冲区数据结构待定
	Compressor Compressor
	Codec      Codec
	buffer     []byte
}
