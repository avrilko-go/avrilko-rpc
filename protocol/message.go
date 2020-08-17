package protocol

import (
	"avrilko-rpc/log"
	"errors"
	"fmt"
	"io"
)

const (
	Magic = 0x08
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

	_, err = io.ReadFull(rBuff, m.Header[1:])
	if err != nil {
		return err
	}



}

var zeroHeaderArr Header
var zeroHeader = zeroHeaderArr[1:]

// 重置头部数据
func resetHeader(m *Header) {
	copy(m[1:], zeroHeader)
}
