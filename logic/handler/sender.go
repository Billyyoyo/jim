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
	clients map[string]rpc.SocketServiceClient
	lock    sync.Mutex
)

func init() {
	// todo 从zookeeper获取所有socket接入服务器
	clients = map[string]rpc.SocketServiceClient{}
	serverUrl := "localhost:5000"
	cli, err := createClient(serverUrl)
	if err != nil {
		return
	}
	clients[serverUrl] = cli
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
	cli := clients[connServer]
	if cli != nil {
		cli.SendAction(context.Background(), action)
	}
}

func SendMessage(connServer string, message *rpc.Message) {
	cli := clients[connServer]
	fmt.Println("server:", connServer, "send message", message.Type.String(), "to", message.DeviceId)
	if cli != nil {
		cli.SendMessage(context.Background(), message)
	}
}

func SendNotification(connServer string, notification *rpc.Notification) {
	fmt.Println("server:", connServer, "send notification", notification.Content, "to", notification.DeviceId)
	cli := clients[connServer]
	if cli != nil {
		cli.SendNotification(context.Background(), notification)
	}
}

func SendKickoff(connServer string, addr *rpc.Text) {
	fmt.Println("server:", connServer, "kickoff", addr.Value)
	cli := clients[connServer]
	if cli != nil {
		cli.SendKickoff(context.Background(), addr)
	}
}
