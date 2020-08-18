package protocol

type Compress interface {
	Zip([]byte) ([]byte, error)
	Unzip([]byte) ([]byte, error)
}
type GzipCompress struct {
}

func (g *GzipCompress) Zip(data []byte) ([]byte, error) {
	panic("implement me")
}

func (g *GzipCompress) Unzip(data []byte) ([]byte, error) {
	panic("implement me")
}
