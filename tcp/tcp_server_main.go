package main

import (
	"jim/common/tool"
	"jim/tcp/server"
	"jim/tcp/service"
)



func main() {
	defer tool.ReleaseGoPool()
	// 启动rpc服务用于接收logic服务器来的消息
	go service.StartUpRpcService()
	// 启动socket服务器
	server.StartUpSocketServer()
}
