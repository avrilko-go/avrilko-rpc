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
	d := []byte(a)
	fmt.Println(d)

	fmt.Println(util.SliceByteToString(d))

}
