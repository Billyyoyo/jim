package tests

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"jim/common/rpc"
	"testing"
)

var (
	cli rpc.LogicServiceClient
)

func init() {
	conn, err := grpc.Dial("localhost:4000", grpc.WithInsecure())
	if err != nil {
		panic("grpc start up error: " + err.Error())
		return
	}
	cli = rpc.NewLogicServiceClient(conn)
}

func TestRpcCreateSession(t *testing.T) {
	req := &rpc.CreateSessionReq{
		Name:    "Test Create Session",
		Creater: 1,
		UserIds: []int64{1, 2, 3},
	}
	session, err := cli.CreateSession(context.Background(), req)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	print(session)
}
