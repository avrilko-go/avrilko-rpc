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
	b := []byte(a)
	fmt.Println(b)
	c, err := util.Zip(b)
	if err != nil {
		panic(err)
	}
	fmt.Println(c)

	d, err := util.Unzip(c)
	if err != nil {
		panic(err)
	}

	fmt.Println(d)

}
