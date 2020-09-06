package sophonn

type Codec interface {
	Read([]byte) []byte
}

type CodecDemo struct {
}

func (c *CodecDemo) Read(data []byte) []byte {
	return data
}
