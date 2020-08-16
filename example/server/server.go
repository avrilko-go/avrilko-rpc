package main

import (
	"avrilko-rpc/example"
	"avrilko-rpc/server"
)

func main() {
	s := server.NewServer()
	he := &example.Hello{}
	s.Register(he, "")
}
