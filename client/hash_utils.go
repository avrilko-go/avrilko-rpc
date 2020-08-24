package client

import (
	"fmt"
	"hash/fnv"
)

func HashString(s string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(s))
	return h.Sum64()
}

// 生成key
func genKey(options ...interface{}) uint64 {
	keyString := ""
	for _, opt := range options {
		keyString = keyString + "/" + toString(opt)
	}

	return HashString(keyString)
}

// 将任何对象转化为字符串
func toString(obj interface{}) string {
	return fmt.Sprintf("%v", obj)
}
