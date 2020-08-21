package main

import (
	"avrilko-rpc/example"
	"avrilko-rpc/server"
	"context"
)

func main() {
	s := server.NewServer()

	s.RegisterFunc(Hb, "")
	err := s.Serve("tcp", "0.0.0.0:8888")
	if err != nil {
		panic(err)
	}
}

func Hb(ctx context.Context, request *example.Request, response *example.Response) error {
	response.C = request.A + request.B
	return nil
}
