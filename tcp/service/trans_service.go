package service

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"jim/common/rpc"
	"jim/tcp/core"
	"jim/tcp/server"
	"net"
)

type TransService struct {
}

func (s *TransService) SendMessage(ctx context.Context, req *rpc.Message) (code *rpc.Int32, err error) {
	log.Info("send msg to ", req.RemoteAddr)
	conn := server.GetUserConn(req.RemoteAddr)
	if conn == nil {
		code = &rpc.Int32{Value: 1}
		return
	}
	body, err := proto.Marshal(req)
	if err != nil {
		return
	}
	pack := rpc.Output{
		Type: rpc.PackType_PT_MSG,
		Data: body,
	}
	bs, err := proto.Marshal(&pack)
	if err != nil {
		return
	}
	(*conn).AsyncWrite(bs)
	code = &rpc.Int32{Value: 0}
	return
}

func (s *TransService) SendNotification(ctx context.Context, req *rpc.Notification) (empty *rpc.Int32, err error) {
	return
}

func (s *TransService) SendAction(ctx context.Context, req *rpc.Action) (code *rpc.Int32, err error) {
	conn := server.GetUserConn(req.RemoteAddr)
	if conn == nil {
		code = &rpc.Int32{Value: 1}
		return
	}
	body, err := proto.Marshal(req)
	if err != nil {
		return
	}
	pack := rpc.Output{
		Type: rpc.PackType_PT_ACTION,
		Data: body,
	}
	bs, err := proto.Marshal(&pack)
	if err != nil {
		return
	}
	(*conn).AsyncWrite(bs)
	code = &rpc.Int32{Value: 0}
	return
}

func (s *TransService) SendKickoff(ctx context.Context, addr *rpc.Text) (empty *rpc.Empty, err error) {
	log.Info("user online duplicate, so kickoff the old connection")
	conn := server.GetUserConn(addr.Value)
	(*conn).Close()
	empty = &rpc.Empty{}
	return
}

func StartUpRpcService() {
	log.Info("=====init grpc server=====")
	addr := fmt.Sprintf("%s:%d", core.AppConfig.Rpc.Host, core.AppConfig.Rpc.Port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
		return
	}
	rpcSever := grpc.NewServer()
	rpc.RegisterSocketServiceServer(rpcSever, &TransService{})
	log.Info("start listen ", addr)
	err = rpcSever.Serve(listener)
	if err != nil {
		panic(err)
	}
}
