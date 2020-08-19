package main

import (
	"avrilko-rpc/util"
	"fmt"
)

type Test struct {
	a []byte
}

func main() {
	a := "我是一个字符串"
	fmt.Println(util.StringToByteSlice(a))
	fmt.Println([]byte(a))

}
