package protocol

import "sync"

var msgPool = &sync.Pool{
	New: func() interface{} {
		header := Header([12]byte{})
		header[0] = Magic
		return &Message{
			Header: &header,
		}
	},
}

// 从缓存池里拿对象
func GetPooledMsg() *Message {
	return msgPool.Get().(*Message)
}

// 重置消息
func FreeMsg(m *Message) {
	if m != nil {
		m.OReset()
		msgPool.Put(m)
	}
}

// 4字节对象池
var poolUint32Data = &sync.Pool{
	New: func() interface{} {
		data := make([]byte, 4)
		return &data
	},
}
