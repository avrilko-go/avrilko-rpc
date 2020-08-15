package main

import (
	"fmt"
	"sync/atomic"
)

func main() {
	var a int32 = 0
	atomic.CompareAndSwapInt32(&a, 1, 0)
	fmt.Println(a)
}
