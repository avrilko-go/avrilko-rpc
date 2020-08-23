package client

type FailMode int

const (
	Failover   FailMode = iota // 自动选择另外一个服务
	Failfast                   // 直接返回错误
	Failtry                    // 使用当前的服务重试
	Failbackup                 //选择另外一个服务器，看看谁返回的快就用谁的返回数据
)

type SelectMode int

const (
	RandomSelect       SelectMode = iota // 随机算法
	RoundRobin                           // 轮询
	WeightedRoundRobin                   // 加权轮询
	WeightedICMP                         // 根据ping值的来加权
	ConsistentHash                       // 一致性hash
	Closest                              // 选择最近的服务器

	SelectByUser = 1000 // 这个由用户自己指定
)


