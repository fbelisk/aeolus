package sophonn

type Codec interface {
	 Read([]byte) ([]byte, error)
}

