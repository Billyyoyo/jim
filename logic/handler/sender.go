package handler

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"jim/common/rpc"
	"sync"
)

var (
	//clients map[string]rpc.SocketServiceClient
	clients sync.Map
)

func init() {
	// todo 从zookeeper获取所有socket接入服务器
	//clients = map[string]rpc.SocketServiceClient{}
	serverUrl := "localhost:4003"
	cli, err := createClient(serverUrl)
	if err != nil {
		return
	}
	clients.Store(serverUrl, &cli)
}

// 给每一个socket服务器创建一个grpc客户端
func createClient(server string) (cli rpc.SocketServiceClient, err error) {
	conn, err := grpc.Dial(server, grpc.WithInsecure())
	if err != nil {
		log.Error("create socket client fail: " + err.Error())
		return
	}
	cli = rpc.NewSocketServiceClient(conn)
	return
}

func SendAction(connServer string, action *rpc.Action) {
	fmt.Println("server:", connServer, "send action", action.Type.String(), "to", action.UserId)
	obj, ok := clients.Load("connServer")
	if ok {
		cli := obj.(rpc.SocketServiceClient)
		cli.SendAction(context.Background(), action)
	}
}

func SendMessage(connServer string, message *rpc.Message) {
	fmt.Println("server:", connServer, "send message", message.Type.String(), "to", message.RemoteAddr)
	obj, ok := clients.Load("connServer")
	if ok {
		cli := obj.(rpc.SocketServiceClient)
		cli.SendMessage(context.Background(), message)
	}
}

func SendNotification(connServer string, notification *rpc.Notification) {
	fmt.Println("server:", connServer, "send notification", notification.Content, "to", notification.DeviceId)
	obj, ok := clients.Load("connServer")
	if ok {
		cli := obj.(rpc.SocketServiceClient)
		cli.SendNotification(context.Background(), notification)
	}
}

func SendKickoff(connServer string, addr *rpc.Text) {
	fmt.Println("server:", connServer, "kickoff", addr.Value)
	obj, ok := clients.Load("connServer")
	if ok {
		cli := obj.(rpc.SocketServiceClient)
		cli.SendKickoff(context.Background(), addr)
	}
}
