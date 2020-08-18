package protocol

import "avrilko-rpc/util"

type Compress interface {
	Zip([]byte) ([]byte, error)
	Unzip([]byte) ([]byte, error)
}
type GzipCompress struct {
}

func (g *GzipCompress) Zip(data []byte) ([]byte, error) {
	return util.Zip(data)
}

func (g *GzipCompress) Unzip(data []byte) ([]byte, error) {
	return util.Unzip(data)
}
