package core

import (
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
	"jim/common/rpc"
	"strconv"
	"time"
)

func (cli *IMClient) handleAuth(pack *rpc.Output) {
	if pack.Code > 0 {
		err := errors.New(pack.Info)
		panic(err)
		return
	} else {
		log.Info("Register device success")
		didr := &rpc.Int64{}
		proto.Unmarshal(pack.Data, didr)
		cli.Ctx.DeviceId = didr.Value
		cli.Status = true
		go cli.heartBeat()
		cli.SyncMsg()
	}
}

func (cli *IMClient) handlePing() {
	pongPack := &rpc.Input{Type: rpc.PackType_PT_PONG,}
	bs, err := proto.Marshal(pongPack)
	if err != nil {
		return
	}
	cli.send(bs)
}

func (cli *IMClient) handlePong() {
	cli.Ctx.PongTime = time.Now().Unix()
}

func (cli *IMClient) handleMsg(msg *rpc.Message) {
	words := &rpc.Words{}
	err := proto.Unmarshal(msg.Content, words)
	if err != nil {
		log.Error("receive a wrong msg, ", err.Error())
		return
	}
	str := fmt.Sprintf("(%d)%d say: %s", msg.SessionId, msg.SendId, words.Text)
	fmt.Println(str)
	ack := rpc.Ack{
		ObjId: msg.Id,
		Type:  rpc.AckType_AT_MESSAGE,
		Seq:   msg.SequenceNo,
	}
	ackbs, err := proto.Marshal(&ack)
	if err != nil {
		log.Error("cant serialliaze ack", err.Error())
		return
	}
	cli.sendPack(rpc.PackType_PT_ACK, &ackbs)
}

func (cli *IMClient) handleAction(act *rpc.Action) {
	str := fmt.Sprintf("%d execute %s", act.UserId, act.Type.String())
	fmt.Println(str)
}

func (cli *IMClient) handleAck(code int32, ack *rpc.Ack) {
	reqId := strconv.FormatInt(ack.RequestId, 10)
	if ack.Type == rpc.AckType_AT_ACT {
		if code > 0 {
			fmt.Println("requestId: " + reqId + " action executed failed.")
		} else {
			fmt.Println("requestId: " + reqId + " action executed success.")
		}
	} else if ack.Type == rpc.AckType_AT_MESSAGE {
		if code > 0 {
			fmt.Println("requestId: " + reqId + " msg sent failed.")
		} else {
			fmt.Println("requestId: " + reqId + " msg sent success.")
		}
	}
}
