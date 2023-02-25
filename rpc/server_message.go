package rpc

import (
	"context"

	service "github.com/indiependente/gw-example/rpc/service/v1"
)

func (s *MsgAPISrv) Echo(ctx context.Context, msg *service.EchoRequest) (*service.EchoResponse, error) {
	s.Log.Println("Received msg: " + msg.Value)
	return &service.EchoResponse{Value: msg.Value}, nil
}
func (s *MsgAPISrv) Reverse(ctx context.Context, msg *service.ReverseRequest) (*service.ReverseResponse, error) {
	s.Log.Println("Received msg: " + msg.Value)
	return &service.ReverseResponse{Value: reverse(msg.Value)}, nil
}

func reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
