package main

import "fmt"

type Test struct {
	a map[string]string
}

func main() {
	a := &Test{
		a: make(map[string]string),
	}

	c := a.a

	c["hb"] = "avrilko"

	fmt.Println(a.a)

}
