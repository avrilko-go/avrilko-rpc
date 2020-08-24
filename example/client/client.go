package main

import (
	"fmt"
	"time"
)

func main() {
	p := make(chan []int)
	go hb(p)

	a := make([]int, 0)
	for i := 0; i < 10; i++ {
		a = append(a, i)
	}

	p <- a

	b := make([]int, 0)
	for j := 0; j < 10; j++ {
		b = append(b, j)
	}

	time.Sleep(5 * time.Second)
	p <- b
	time.Sleep(5 * time.Second)
}

func hb(ch <-chan []int) {
	for v := range ch {
		for _, vv := range v {
			fmt.Println(vv)
		}
	}

	fmt.Println(111)
}
