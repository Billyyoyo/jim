package main

import (
	"bufio"
	log "github.com/sirupsen/logrus"
	"jim/client/core"
	"os"
)

var (
	client *core.IMClient
)

func main() {
	// todo 1.请求http服务器获取连接地址
	// todo 2.链接接入服务器
	// todo 3.发送设备注册指令
	// todo 4.同步消息
	addr := "localhost:4002"
	var uid int64 = 1
	var err error
	client, err = core.NewIMClient(uid, addr)
	if err != nil {
		panic(err.Error())
		return
	}
	log.Info("Connected to ", addr, " success")
	inputer()
}

func inputer() {
	reader := bufio.NewReader(os.Stdin)
	for {
		str, _ := reader.ReadString('\n')
		client.Command([]byte(str))
	}
}
