package util

import (
	"math"
	"sync"
)

type levelPool struct {
	size int
	pool *sync.Pool
}

type LimitedPool struct {
	minSize int
	maxSize int
	pools   []*levelPool
}

// 实例化一个固定大小的对象缓存池子
func newLevelPool(size int) *levelPool {
	return &levelPool{
		size: size,
		pool: &sync.Pool{
			New: func() interface{} {
				data := make([]byte, size)
				return &data
			},
		},
	}
}

// new一个指定范围的链接池子
func NewLimitedPool(minSize, maxSize int) *LimitedPool {
	if maxSize < minSize {
		return &LimitedPool{}
	}

	const StretchSize = 2 // 伸缩比
	curSize := minSize
	var pools []*levelPool

	for curSize < maxSize { //2的倍数
		pools = append(pools, newLevelPool(curSize))
		curSize *= StretchSize
	}

	pools = append(pools, newLevelPool(maxSize))

	return &LimitedPool{
		minSize: minSize,
		maxSize: maxSize,
		pools:   pools,
	}
}

func (l *LimitedPool) Get(size int) *[]byte {
	levelPool := l.findPool(size)
	if levelPool == nil {
		data := make([]byte, size)
		return &data
	}

	data := levelPool.pool.Get().(*[]byte)
	*data = (*data)[:size] // 只取有效部分
	return data
}

// 查找指定范围的对象缓存池
func (l *LimitedPool) findPool(size int) *levelPool {
	if size > l.maxSize {
		return nil
	}

	index := int(math.Ceil(math.Log2(float64(size) / float64(l.minSize))))
	if index < 0 {
		index = 0
	}

	if index > (len(l.pools) - 1) {
		return nil
	}

	return l.pools[index]
}
