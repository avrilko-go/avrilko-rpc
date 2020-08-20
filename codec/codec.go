package codec

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gogo/protobuf/proto"
	pb "github.com/golang/protobuf/proto"
	"github.com/vmihailenco/msgpack"
	"reflect"
)

type Codec interface {
	Encode(interface{}) ([]byte, error) // 序列化
	Decode([]byte, interface{}) error   // 反序列化
}

// 使用原始的byte切片来传输，要求传入的是[]byte或者*[]byte类型
type ByteCodec struct {
}

func (b *ByteCodec) Encode(i interface{}) ([]byte, error) {
	if data, ok := i.([]byte); ok {
		return data, nil
	}
	if data, ok := i.(*[]byte); ok {
		return *data, nil
	}

	return nil, errors.New(fmt.Sprintf("序列化失败，传入的类型%T不是[]byte类型或者*[]byte类型", i))
}

func (b *ByteCodec) Decode(bytes []byte, i interface{}) error {
	reflect.Indirect(reflect.ValueOf(i)).SetBytes(bytes)
	return nil
}

type JSONCodec struct {
}

func (J *JSONCodec) Encode(i interface{}) ([]byte, error) {
	return json.Marshal(i)
}

func (J *JSONCodec) Decode(bytes []byte, i interface{}) error {
	return json.Unmarshal(bytes, i)
}

// proto 兼容两种方式
type PBCodec struct {
}

func (P *PBCodec) Encode(i interface{}) ([]byte, error) {
	if m, ok := i.(proto.Marshaler); ok {
		return m.Marshal()
	}
	if m, ok := i.(pb.Message); ok {
		return pb.Marshal(m)
	}

	return nil, fmt.Errorf("传入的参数%T不是一个proto.Marshaler", i)

}

func (P *PBCodec) Decode(bytes []byte, i interface{}) error {
	if m, ok := i.(proto.Unmarshaler); ok {
		return m.Unmarshal(bytes)
	}
	if m, ok := i.(pb.Message); ok {
		return pb.Unmarshal(bytes, m)
	}

	return fmt.Errorf("传入的参数%T不是一个proto.UnMarshaler", i)
}

// msgpack 方式序列化  默认支持json tag
type MsgpackCodec struct {
}

func (m *MsgpackCodec) Encode(i interface{}) ([]byte, error) {
	buff := bytes.NewBuffer(nil)
	en := msgpack.NewEncoder(buff)
	en.UseJSONTag(true)
	err := en.Encode(i)
	return buff.Bytes(), err
}

func (m *MsgpackCodec) Decode(data []byte, i interface{}) error {
	b := bytes.NewBuffer(data)
	de := msgpack.NewDecoder(b)
	de.UseJSONTag(true)
	return de.Decode(i)
}
