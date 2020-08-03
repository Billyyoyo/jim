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
	clients.Store(serverUrl, cli)
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

func SendAction(connServer string, action *rpc.Action) (ret int32) {
	fmt.Println("server:", connServer, "send action", action.Type.String(), "to", action.UserId)
	obj, ok := clients.Load(connServer)
	if ok {
		cli := obj.(rpc.SocketServiceClient)
		code, err := cli.SendAction(context.Background(), action)
		if err != nil {
			log.Error("send action to tcp fail: ", err.Error())
			return
		}
		ret = code.Value
	}
	return
}

func SendMessage(connServer string, message *rpc.Message)(ret int32) {
	fmt.Println("server:", connServer, "send message", message.Type.String(), "to", message.RemoteAddr)
	obj, ok := clients.Load(connServer)
	if ok {
		cli := obj.(rpc.SocketServiceClient)
		code, err := cli.SendMessage(context.Background(), message)
		if err != nil {
			log.Error("send message to tcp fail: ", err.Error())
		}
		ret = code.Value
	}
	return
}

func SendNotification(connServer string, notification *rpc.Notification)(ret int32) {
	fmt.Println("server:", connServer, "send notification", notification.Content, "to", notification.DeviceId)
	obj, ok := clients.Load(connServer)
	if ok {
		cli := obj.(rpc.SocketServiceClient)
		code, err := cli.SendNotification(context.Background(), notification)
		if err != nil {
			log.Error("send notify to tcp fail: ", err.Error())
		}
		ret = code.Value
	}
	return
}

func SendKickoff(connServer string, addr *rpc.Text) {
	fmt.Println("server:", connServer, "kickoff", addr.Value)
	obj, ok := clients.Load(connServer)
	if ok {
		cli := obj.(rpc.SocketServiceClient)
		_, err := cli.SendKickoff(context.Background(), addr)
		if err != nil {
			log.Error("kick off someone to tcp fail: ", err.Error())
		}
	}
}
