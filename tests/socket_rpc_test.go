package tests

import (
	"context"
	"github.com/golang/protobuf/proto"
	"jim/common/rpc"
	"testing"
)

func TestSocketRpcSendMessage(t *testing.T) {
	words:=&rpc.Words{
		Text:                 "sfasdf",
	}
	bs,_:=proto.Marshal(words)
	msg := &rpc.Message{
		Id:                   1,
		SendId:               1,
		SessionId:            9,
		Time:                 0,
		RequestId:            0,
		Status:               1,
		Type:                 0,
		Content:              bs,
		SequenceNo:           0,
		RemoteAddr:           "127.0.0.1:43514",
		XXX_NoUnkeyedLiteral: struct{}{},
		XXX_unrecognized:     nil,
		XXX_sizecache:        0,
	}
	_, err := cli2.SendMessage(context.Background(), msg)
	if err != nil {
		printl(err.Error())
	}
}
