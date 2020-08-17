package main

import (
	"fmt"
	"time"
)

func main() {

	now := time.Now()

	fmt.Println(now.Add(time.Second * 5))
	fmt.Println(now.Add(time.Second * 5))

}
