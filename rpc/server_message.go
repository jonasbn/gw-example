package rpc

import (
	"context"
)

func (s *MsgAPISrv) Echo(ctx context.Context, msg *StringMessage) (*StringMessage, error) {
	s.Log.Println("Received msg: " + msg.Value)
	return msg, nil
}
func (s *MsgAPISrv) Reverse(ctx context.Context, msg *StringMessage) (*StringMessage, error) {
	s.Log.Println("Received msg: " + msg.Value)
	return &StringMessage{Value: reverse(msg.Value)}, nil
}

func reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
