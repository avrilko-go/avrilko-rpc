package protocol

import (
	"avrilko-rpc/log"
	"avrilko-rpc/util"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

const (
	ServiceError = "__rpcx_error__"
)

const (
	Magic = 0x08
)

var (
	CompressTypeMap = map[CompressType]Compress{
		Gzip: &GzipCompress{},
	}
)

var (
	MaxMessageLength = 0 // 最大消息体长度（不包括头部）为0则无限制
)

var (
	ErrMessageTooLong      = errors.New("消息体的长度太长了，超过最大长度")
	ErrMetaKVMissing       = errors.New("错误的meta信息，可能丢失了数据")
	ErrUnsupportedCompress = errors.New("不支持的压缩类型")
)

// 数据压缩的类型
type CompressType byte

const (
	None CompressType = iota // 不压缩
	Gzip                     // 使用gzip压缩
)

// 消息的方式
type MessageType byte

const (
	Request MessageType = iota
	Response
)

// 消息类型
type MessageStatusType byte

const (
	Normal MessageStatusType = iota
	Error
)

// 获取框架的魔数
func MagicNumber() byte {
	return Magic
}

// 头部包括4字节的Header + 8字节的Message(seq)
type Header [12]byte

func (h *Header) CheckMagic() bool {
	return h[0] == Magic
}

func (h Header) MessageStatusType() MessageStatusType {
	return MessageStatusType(h[2] & 0x03)
}

func (h *Header) SetMessageStatusType(m MessageStatusType) {
	h[2] = (h[2] &^ 0x03) | (byte(m) & 0x03)
}

func (h *Header) SetMessageType(m MessageType) {
	h[2] = (h[2] &^ 0x80) | (byte(m) << 7)
}

func (h Header) MessageType() MessageType {
	return MessageType(h[2]&0x80) >> 7
}

func (h Header) IsHeartbeat() bool {
	return h[2]&0x40 == 0x40
}

func (h *Header) SetHeartbeat(hb bool) {
	if hb {
		h[2] = h[2] | 0x40
	} else { // 这是个清零操作， 只要对比的位上为1则被对比的位上就会清0
		h[2] = h[2] &^ 0x40
	}
}

func (h Header) IsOneway() bool {
	return h[2]&0x20 == 0x20
}

func (h *Header) SetOneway(one bool) {
	if one {
		h[2] = h[2] | 0x20
	} else {
		h[2] = h[2] &^ 0x20
	}
}

// 获取压缩类型 (***xxx** & 00011100) => (xxx00) >> 2 => (xxx) （蛋疼的算法）
func (h Header) CompressType() CompressType {
	// 这里不用引用（使用值拷贝）
	return CompressType(h[2]&0x1c) >> 2
}

// 设置压缩类型
// 首先使用&^ 符号初始化指定的3位，将其都变成000 然后或操作将值赋予指定的3位
func (h *Header) SetCompressType(c CompressType) {
	h[2] = (h[2] &^ 0x1c) | (byte(c) << 2)
}

// rpc 标准的请求和响应格式
type Message struct {
	*Header                         // 头部信息（包括魔数 + 版本号 + 消息类型 + 是否是心跳 + 是否是上报服务 + 是否压缩 + 单个请求是否是成功 + 序列化方式）
	ServicePath   string            // 路由地址
	ServiceMethod string            // 路由方法
	Metadata      map[string]string // 元数据（穿透服务端和客户端的，可用来做鉴权）
	Payload       []byte            // 真正传输的数据（客户端和服务端都放在这）
	data          []byte            // 工具人 除了头部和整个数据长度以外的其他数据(有点工具人的感觉)
}

// 重置消息体
func (m *Message) OReset() {
	resetHeader(m.Header)
	m.ServiceMethod = ""
	m.ServicePath = ""
	m.data = []byte{}
	m.Payload = []byte{}
	m.Metadata = nil
}

// 解码数据
func (m *Message) Decode(rBuff io.Reader) (err error) {
	_, err = io.ReadFull(rBuff, m.Header[:1]) // 只读一个字节
	if err != nil {
		return err
	}

	if !m.Header.CheckMagic() {
		errStr := fmt.Sprintf("不是avrilko-rpc协议 魔数为%d", m.Header[0])
		log.Error(errStr)
		return errors.New(errStr)
	}

	_, err = io.ReadFull(rBuff, m.Header[1:]) // 将头部全部读出来
	if err != nil {
		return err
	}

	totalLen := poolUint32Data.Get().(*[]byte) // 取出来是指向[]byte的指针
	_, err = io.ReadFull(rBuff, *totalLen)     // 这里用*是取其值的意思
	if err != nil {
		poolUint32Data.Put(totalLen)
		return err
	}
	l := binary.BigEndian.Uint32(*totalLen)
	poolUint32Data.Put(totalLen)

	totalDataLen := int(l)
	if MaxMessageLength > 0 && totalDataLen > MaxMessageLength {
		return ErrMessageTooLong
	}

	if cap(m.data) > totalDataLen { // 自身容量大了缩容
		m.data = m.data[:totalDataLen]
	} else { // 小了需要开辟空间
		m.data = make([]byte, totalDataLen)
	}

	data := m.data //这里将data单独拿出来，操作data相当于操作m.data

	_, err = io.ReadFull(rBuff, data) // 将所有内容读到data中
	if err != nil {
		return err
	}
	start := 0
	fieldLen := int(binary.BigEndian.Uint32(data[start:4])) // 读出servicePath的长度
	start += 4
	m.ServicePath = util.SliceByteToString(data[start : start+fieldLen])
	start += fieldLen

	fieldLen = int(binary.BigEndian.Uint32(data[start:4])) // 读出serviceMethod的长度
	start += 4
	m.ServiceMethod = util.SliceByteToString(data[start : start+fieldLen])
	start += fieldLen

	fieldLen = int(binary.BigEndian.Uint32(data[start:4])) // 读出meta的长度
	start += 4
	if fieldLen > 0 { // 传递了meta信息则解析
		m.Metadata, err = decodeMetaData(fieldLen, data[start:start+fieldLen])
		if err != nil {
			return err
		}
	}
	start += fieldLen

	fieldLen = int(binary.BigEndian.Uint32(data[start:4])) // 读出payload长度
	start += 4
	_ = start // 销毁变量
	m.Payload = data[start:]
	// 剩下的data数据全部为payload的
	if m.CompressType() != None { // 使用了gzip压缩
		compressImpl, ok := CompressTypeMap[m.CompressType()]
		if !ok {
			return ErrUnsupportedCompress
		}
		m.Payload, err = compressImpl.Unzip(m.Payload)
		if err != nil {
			return err
		}
	}

	return nil
}

// 拷贝一个message对象
func (m Message) Clone() *Message {
	header := *m.Header
	c := GetPooledMsg()
	header.SetCompressType(None)
	c.Header = &header
	c.ServicePath = m.ServicePath
	c.ServiceMethod = m.ServiceMethod
	return c
}

func (m *Message) EncodeSlicePointer() *[]byte {
	// 打包头部
	meta := encodeMetaData(m.Metadata)
	pLen := len(m.ServicePath)
	mLen := len(m.ServiceMethod)

	var err error // 申明一个变量
	payload := m.Payload
	if m.CompressType() != None {
		compress := CompressTypeMap[m.CompressType()]
		if compress == nil {
			m.SetCompressType(None)
		} else {
			payload, err = compress.Zip(m.Payload)
			if err != nil {
				payload = m.Payload
			}
		}
	}

	totalLen := (pLen + 4) + (mLen + 4) + (4 + len(meta)) + (4 + len(payload))
	metaStart := 12 + 4 + pLen + 4 + mLen + 4
	payLoadStart := metaStart + 4 + len(meta)

	l := 12 + 4 + totalLen






}

// 打包meta信息
func encodeMetaData(m map[string]string) []byte {
	if len(m) == 0 {
		return []byte{}
	}

	buff := bytes.NewBuffer(nil)
	lens := make([]byte, 4)
	for k, v := range m {
		// 写入key
		binary.BigEndian.PutUint32(lens, uint32(len(k)))
		buff.Write(lens)
		buff.Write(util.StringToByteSlice(k))

		// 写入value
		binary.BigEndian.PutUint32(lens, uint32(len(v)))
		buff.Write(lens)
		buff.Write(util.StringToByteSlice(v))
	}

	return buff.Bytes()
}

// 解析meta信息
func decodeMetaData(l int, data []byte) (map[string]string, error) {
	m := make(map[string]string, 10)
	n := 0
	lenKV := 0
	key := ""
	value := ""
	for n < l {
		// 每个key和value的长度都是4个字节的字符串
		lenKV = int(binary.BigEndian.Uint32(data[n:4]))
		n += 4
		if n+lenKV < l-4 { // 没法解析了
			return m, ErrMetaKVMissing
		}
		key = util.SliceByteToString(data[n : n+lenKV])
		n += lenKV

		lenKV = int(binary.BigEndian.Uint32(data[n:4]))
		n += 4
		if n < l-4 { // 没法解析了
			return m, ErrMetaKVMissing
		}
		value = util.SliceByteToString(data[n : n+lenKV])
		n += lenKV

		m[key] = value
	}

	return m, nil
}

var zeroHeaderArr Header
var zeroHeader = zeroHeaderArr[1:]

// 重置头部数据
func resetHeader(m *Header) {
	copy(m[1:], zeroHeader)
}
