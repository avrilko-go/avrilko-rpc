package main

import (
	"avrilko-rpc/example"
	"avrilko-rpc/server"
	"context"
)

func main() {
	s := server.NewServer()
	he := &example.Hello{}
	s.Register(he, "")
	s.Shutdown(context.TODO())
}
