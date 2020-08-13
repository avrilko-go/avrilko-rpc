package serilazation

// 包装好的解包函数，在注册服务时就注册进去，具体怎么解包由调用者传入
type WrapperUnmarshal func(data interface{}) error
