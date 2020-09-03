package sophonn

type Compressor interface {
	Decode([]byte)[]byte
	Encode([]byte)[]byte
}

