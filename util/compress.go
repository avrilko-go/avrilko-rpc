package util

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
	"sync"
)

var (
	spReader *sync.Pool
	spWriter *sync.Pool
	spBuffer *sync.Pool
)

func init() {
	spReader = &sync.Pool{
		New: func() interface{} {
			return new(gzip.Reader)
		},
	}
	spWriter = &sync.Pool{
		New: func() interface{} {
			return new(gzip.Writer)
		},
	}
	spBuffer = &sync.Pool{
		New: func() interface{} {
			return bytes.NewBuffer(nil)
		},
	}
}

// 使用gzip压缩
func Zip(data []byte) ([]byte, error) {
	buff := spBuffer.Get().(*bytes.Buffer)
	w := spWriter.Get().(*gzip.Writer)

	w.Reset(buff)
	defer func() {
		buff.Reset()
		spBuffer.Put(buff)
		w.Close()
		spWriter.Put(w)
	}()

	_, err := w.Write(data)
	if err != nil {
		return nil, err
	}
	err = w.Flush()
	if err != nil {
		return nil, err
	}
	err = w.Close()
	if err != nil {
		return nil, err
	}

	return buff.Bytes(), nil
}

// 使用gzip解压缩
func Unzip(data []byte) ([]byte, error) {
	buff := spBuffer.Get().(*bytes.Buffer)
	defer func() {
		buff.Reset()
		spBuffer.Put(buff)
	}()

	_, err := buff.Write(data)
	if err != nil {
		return nil, err
	}

	r := spReader.Get().(*gzip.Reader)
	defer func() {
		spReader.Put(r)
	}()
	err = r.Reset(buff)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	originData, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return originData, nil
}
