package sophonn

type Compressor interface {
	Decode([]byte)[]byte
	Encode([]byte)[]byte
}

type CompressorDemo struct {

}

func (c *CompressorDemo) Decode(data []byte) []byte {
	return data
}

func (c *CompressorDemo) Encode(data []byte) []byte {
	return data
}


