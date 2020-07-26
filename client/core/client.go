package core

import (
	"bufio"
	"errors"
	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
	"io"
	"jim/common/rpc"
	"net"
	"os"
	"time"
)

type IMClient struct {
	conn net.Conn
	Addr string
	Ctx  Ctx
}

type Ctx struct {
	UserId   int64
	DeviceId int64
	PongTime int64
}

func NewIMClient(uid int64, addr string) (client *IMClient, err error) {
	client = &IMClient{
		Addr: addr,
		Ctx: Ctx{
			UserId:   uid,
			PongTime: time.Now().Unix(),
		},
	}
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return
	} else {
		client.conn = conn
		go client.loop()
		go client.heartBeat()
	}
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
		inPack := &rpc.Input{}
		err = proto.Unmarshal(ret, inPack)
		if err != nil {
			log.Error(err.Error())
			continue
		}
		cli.handlePack(inPack)
	}
}

func (cli *IMClient) handlePack(pack *rpc.Input) {
	if pack.Type == rpc.PackType_PT_PING {
		pongPack := &rpc.Output{Type: rpc.PackType_PT_PONG,}
		bs, err := proto.Marshal(pongPack)
		if err != nil {
			return
		}
		cli.send(bs)
	} else if pack.Type == rpc.PackType_PT_PONG {
		cli.Ctx.PongTime = time.Now().Unix()
	}
}

//发送ping到server，收到pong确认连接，否则超时断开连接
func (cli *IMClient) heartBeat() {
	tick := time.Tick(10 * time.Second)
	for {
		select {
		case <-tick:
			if time.Now().Unix()-cli.Ctx.PongTime > 30 {
				// 心跳超时将关闭连接，并退出
				log.Info("heart beat timeout")
				err := cli.conn.Close()
				if err != nil {
					panic(err.Error())
					return
				}
				os.Exit(1)
			}
			pingPack := &rpc.Output{Type: rpc.PackType_PT_PING,}
			bs, err := proto.Marshal(pingPack)
			if err != nil {
				return
			}
			cli.send(bs)
			break
		}
	}
}

// :为动作 如加入会话，退出会话等 :join 10, :rename 10 码农群  :create 码农群 1,2,3
// /为发送消息  如 /10 大家好
// -为撤销消息   -1003
func (cli *IMClient) Command(cmd []byte) (err error) {
	return
}
