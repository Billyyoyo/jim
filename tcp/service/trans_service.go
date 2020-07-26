package service

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"jim/common/rpc"
	"jim/tcp/core"
	"net"
)

type TransService struct {
}

func (s *TransService) SendMessage(ctx context.Context, req *rpc.Message) (empty *rpc.Empty, err error) {
	return
}

func (s *TransService) SendNotification(ctx context.Context, req *rpc.Notification) (empty *rpc.Empty, err error) {
	return
}

func (s *TransService) SendAction(ctx context.Context, req *rpc.Action) (empty *rpc.Empty, err error) {
	return
}

func (s *TransService) SendKickoff(ctx context.Context, addr *rpc.Text) (empty *rpc.Empty, err error) {
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
