package sophonn

import "github.com/fbelisk/aeolus/ringbuffer"

//todo 抽象conn buffer 处理 编解码 压缩 及 读取状态存储的问题
type ConnBuffer struct {
	buffer *ringbuffer.RingBuffer
	Compressor
	Codec
}

func NewConnBuffer(size int) *ConnBuffer {
	r := &ConnBuffer{
		buffer: ringbuffer.New(0x1000),
		Compressor:&CompressorDemo{},
		Codec:&CodecDemo{},
	}
	return r
}