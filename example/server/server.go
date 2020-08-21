package main

import (
	"avrilko-rpc/example"
	"avrilko-rpc/server"
)

func main() {
	s := server.NewServer()
	h := &example.Hello{}

	s.Register(h, "")
	err := s.Serve("tcp", "0.0.0.0:8888")
	if err != nil {
		panic(err)
	}
}
