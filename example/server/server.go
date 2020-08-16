package main

import (
	"avrilko-rpc/example"
	"fmt"
	"reflect"
)



func main() {
	a := &example.Hb{}

	r := reflect.TypeOf(a)
	fmt.Println(r.Kind() == reflect.Ptr)
	fmt.Println(r.Elem().PkgPath())

}
