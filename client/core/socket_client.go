package core

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
	"io"
	"jim/common/rpc"
	"net"
	"os"
	"strings"
	"time"
)

type IMClient struct {
	conn   net.Conn
	Addr   string
	Ctx    Ctx
	Status bool
}

type Ctx struct {
	UserId    int64
	DeviceId  int64
	Token     string
	Serial    string
	PongTime  int64
	SessionId int64
}

func NewIMClient(uid int64, token, serial, socketUrl string) (client *IMClient, err error) {
	client = &IMClient{
		Addr: socketUrl,
		Ctx: Ctx{
			UserId:   uid,
			Token:    token,
			Serial:   serial,
			PongTime: time.Now().Unix(),
		},
	}
	conn, err := net.Dial("tcp", socketUrl)
	if err != nil {
		return
	} else {
		client.conn = conn
		go client.loop()
		go client.reg()
	}
	return
}

func (cli *IMClient) reg() {
	info := &rpc.RegInfo{
		UserId:   cli.Ctx.UserId,
		Token:    cli.Ctx.Token,
		SerialNo: cli.Ctx.Serial,
	}
	bs, err := proto.Marshal(info)
	if err != nil {
		panic(err)
	}
	err = cli.sendPack(rpc.PackType_PT_AUTH, &bs)
	if err != nil {
		panic(err)
	}
}

func (cli *IMClient) sendPack(pt rpc.PackType, bs *[]byte) (err error) {
	pack := &rpc.Input{
		Type: pt,
		Data: *bs,
	}
	bb, err := proto.Marshal(pack)
	if err != nil {
		log.Error("send pack - ", err.Error())
		return
	}
	cli.send(bb)
	return
}

func (cli *IMClient) sendAction(at rpc.ActType, bs *[]byte) (err error) {
	act := &rpc.Action{
		UserId:    cli.Ctx.UserId,
		RequestId: time.Now().Unix(),
		Time:      time.Now().Unix(),
		Type:      at,
		Data:      *bs,
	}
	*bs, err = proto.Marshal(act)
	if err != nil {
		log.Error("send act - ", err.Error())
		return
	}
	cli.sendPack(rpc.PackType_PT_ACTION, bs)
	return
}

func (cli *IMClient) send(data []byte) {
	pack, err := Encode(data)
	if err != nil {
		log.Info("encode error:", err.Error())
		return
	}
	_, err = cli.conn.Write(pack)
	if err != nil {
		log.Info("write error:", err.Error())
	}
}

func (cli *IMClient) loop() {
	for {
		inputReader := bufio.NewReader(cli.conn)
		ret, err := Decode(inputReader)
		if err != nil {
			if io.EOF == err {
				panic(errors.New("connection is closed"))
			}
			continue
		}
		outPack := &rpc.Output{}
		err = proto.Unmarshal(ret, outPack)
		if err != nil {
			log.Error(err.Error())
			continue
		}
		cli.dispatchPack(outPack)
	}
}

func (cli *IMClient) dispatchPack(pack *rpc.Output) {
	if pack.Type == rpc.PackType_PT_PING {
		go cli.handlePing()
	} else if pack.Type == rpc.PackType_PT_PONG {
		go cli.handlePong()
	} else if pack.Type == rpc.PackType_PT_AUTH {
		go cli.handleAuth(pack)
	} else if pack.Type == rpc.PackType_PT_MSG {
		msg := &rpc.Message{}
		if err := proto.Unmarshal(pack.Data, msg); err == nil {
			go cli.handleMsg(msg)
		}
	} else if pack.Type == rpc.PackType_PT_ACTION {
		act := &rpc.Action{}
		if err := proto.Unmarshal(pack.Data, act); err == nil {
			go cli.handleAction(act)
		}
	} else if pack.Type == rpc.PackType_PT_ACK {
		ack := &rpc.Ack{}
		if err := proto.Unmarshal(pack.Data, ack); err == nil {
			go cli.handleAck(pack.Code, pack.Info, ack)
		}
	} else if pack.Type == rpc.PackType_PT_NOTIFICATION {

	}
}

//发送ping到server，收到pong确认连接，否则超时断开连接
func (cli *IMClient) heartBeat() {
	tick := time.Tick(10 * time.Second)
	for {
		select {
		case <-tick:
			if time.Now().Unix()-cli.Ctx.PongTime > 60 {
				// 心跳超时将关闭连接，并退出
				log.Info("heart beat timeout")
				err := cli.conn.Close()
				if err != nil {
					panic(err.Error())
					return
				}
				os.Exit(1)
			}
			pingPack := &rpc.Input{Type: rpc.PackType_PT_PING,}
			bs, err := proto.Marshal(pingPack)
			if err != nil {
				return
			}
			cli.send(bs)
			break
		}
	}
}

// :为动作 如加入会话，退出会话等 :join 10, :rename 10 码农群  :create 码农群 1,2,3 :enter 18
// /为发送消息  如 /10 大家好
// -为撤销消息   -1003
func (cli *IMClient) Command(bs []byte) (err error) {
	if !cli.Status {
		fmt.Println("client not init complete, wait register to im server")
		return
	}
	cmd := string(bs)
	if strings.Index(cmd, ":") == 0 {
		if strings.Index(cmd, " ") < 0 {
			err = errors.New("")
		} else {
			command := cmd[1:strings.Index(cmd, " ")]
			content := cmd[strings.Index(cmd, " ")+1:]
			if command == "ne" {
				cli.CreateSession(content)
				return
			} else if command == "jo" {
				cli.JoinSession(content)
				return
			} else if command == "qu" {
				cli.QuitSession(content)
				return
			} else if command == "re" {
				cli.RenameSession(content)
				return
			} else if command == "-" {
				cli.WithdrawMsg(content)
				return
			} else if command == "sl" {
				cli.GetSessions()
				return
			} else if command == "ml" {
				cli.GetMembers()
				return
			} else if command == "sw" {
				cli.SwitchSession(content)
				return
			} else if command == "ms" {
				cli.GetMessages(content)
				return
			} else if command == "-" {
				return
			} else if command == "sl" {
				return
			} else if command == "ml" {
				return
			} else if command == "sw" {
				return
			} else {
				err = errors.New("")
			}
		}
		if err != nil {
			fmt.Println("no this command. all commands like below:" +
				"\n\t:ne\tcreate session" +
				"\n\t:jo\tjoin session" +
				"\n\t:qu\tquit session" +
				"\n\t:re\trename session" +
				"\n\t:-\twithdraw message" +
				"\n\t:sl mine\tlist all your sessions" +
				"\n\t:ml mine\tlist all members in session" +
				"\n\t:ms\tlist messages in session" +
				"\n\t:sw\tswitch session")
		}
	} else if strings.Index(cmd, "/") == 0 {
		// 发送消息
		cli.sendMsg(cmd[1:])
	} else {
		fmt.Println("Bad command! begin char is: \n\t/xxx send msg \n\t:cmd execute create, join, withdraw or quit command")
	}
	return
}
