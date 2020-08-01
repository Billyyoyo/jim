package main

import (
	"bufio"
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"jim/client/core"
	"os"
	"strings"
)

var (
	authCode string
	serialNo string
)

var (
	client *core.IMClient
)

func init() {
	flag.StringVar(&authCode, "code", "", "授权令牌")
	flag.StringVar(&serialNo, "serial", "", "设备序列号")
}

func main() {
	flag.Parse()
	if authCode == "" {
		fmt.Println("please run client with paramter -code xxxxx")
		return
	}
	if serialNo == "" {
		fmt.Println("please run client with paramter -serial xxxxx")
		return
	}
	// todo 1.请求http服务器获取连接地址
	// todo 2.链接接入服务器
	// todo 3.发送设备注册指令
	// todo 4.同步消息
	signIn()
}

func signIn() {
	uid, token, server, err := core.Authorization(authCode)
	if err != nil {
		panic("request authorization fail:" + err.Error())
		return
	}
	connectServer(uid, token, server)
}

func connectServer(uid int64, token, socketAddr string) {
	var err error
	client, err = core.NewIMClient(uid, token, serialNo, socketAddr)
	if err != nil {
		panic(err.Error())
		return
	}
	log.Info("Connected to ", socketAddr, " success")
	inputer()
}

func inputer() {
	reader := bufio.NewReader(os.Stdin)
	for {
		str, _ := reader.ReadString('\n')
		str = str[0 : len(str)-1]
		str = strings.Trim(str, " ")
		client.Command([]byte(str))
	}
}
