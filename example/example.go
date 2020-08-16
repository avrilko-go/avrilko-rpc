package example

import "context"

type Request struct {
	A int
	B int
}

type Response struct {
	C int
}

type Hello struct {
}

func (h *Hello) Sum(ctx context.Context, request *Request, response *Response) error {
	response.C = request.A + request.B
	return nil
}

//type Hello2 struct {
//}
//
//func (h *Hello2) Sum2(ctx context.Context, request *Request, response *Response) error {
//	response.C = request.A + request.B
//	return nil
//}
//
//func (r *Request) Reset() {
//	r.A = 0
//	r.B = 0
//}
