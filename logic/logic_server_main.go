package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"jim/common/rpc"
	"jim/logic/core"
	_ "jim/logic/core"
	"jim/logic/service"
	"net"
)

func main() {
	log.Info("=====register server on name resolver=====")
	//todo
	log.Info("=====register server on message queue=====")
	//todo
	log.Info("=====init grpc server=====")
	addr := fmt.Sprintf("%s:%d", core.AppConfig.Server.Host, core.AppConfig.Server.Port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
		return
	}
	rpcSever := grpc.NewServer()
	rpc.RegisterLogicServiceServer(rpcSever, &service.LogicService{})
	log.Info("start listen ", addr)
	err = rpcSever.Serve(listener)
	if err != nil {
		panic(err)
	}
}
