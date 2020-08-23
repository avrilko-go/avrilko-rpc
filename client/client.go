package client

type Call struct {
	ServicePath   string
	ServiceMethod string
	Metadata      map[string]string // 需要传递到服务端的元数据
	ResMetadata   map[string]string // 接受服务端返回的元数据
	request       interface{}       // 客户端请求
	response      interface{}       // 响应
	Error         error             // 调用完成之后的错误
	Done          chan *Call        // 调用完整之后会讲数据塞到此通道中
	Raw           bool              // 是否发送原始数据
}
