package tests

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/proto"
	"io"
	"jim/common/rpc"
	"testing"
)

func TestRpcCreateSession(t *testing.T) {
	req := &rpc.CreateSessionReq{
		Name:    "Small group 2",
		Creater: 1,
		Type:    rpc.SessionType_SESSION_GROUP,
		UserIds: []int64{1, 2},
	}
	session, err := cli.CreateSession(context.Background(), req)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	printj(session)
}

func TestRpcJoinSession(t *testing.T) {
	req := &rpc.JoinSessionReq{
		UserId:    4,
		SessionId: 10,
	}
	ret, err := cli.JoinSession(context.Background(), req)
	if err != nil {
		printl(err.Error())
		return
	}
	printj(ret)
}

func TestRpcQuitSession(t *testing.T) {
	req := &rpc.QuitSessionReq{
		UserId:    4,
		SessionId: 10,
	}
	ret, err := cli.QuitSession(context.Background(), req)
	if err != nil {
		printl(err.Error())
		return
	}
	printj(ret)
}

func TestRpcRegister(t *testing.T) {
	req := &rpc.RegisterReq{
		UserId:   1,
		Token:    "749f6530e13a40af92cf17db2c5e2130",
		Addr:     "localhost:23224",
		Server:   "localhost:5000",
		SerialNo: "20139fcd-25fe-42e0-9457-49356018beb8",
	}
	resp, err := cli.Register(context.Background(), req)
	if err != nil {
		printl(err.Error())
		return
	}
	printj(resp)
}

func TestRpcOffline(t *testing.T) {
	req := &rpc.OfflineReq{
		UserId:       1,
		DeviceId:     16,
	}
	_, err := cli.Offline(context.Background(), req)
	if err != nil {
		printl(err.Error())
		return
	}
}

func TestRpcRenameSession(t *testing.T) {
	req := &rpc.RenameSessionReq{
		SessionId: 10,
		UserId:    1,
		Name:      "大家嗨起来",
	}
	ret, err := cli.RenameSession(context.Background(), req)
	if err != nil {
		printl(err.Error())
		return
	}
	printj(ret)
}

func TestRpcGetSession(t *testing.T) {
	sessionId := &rpc.Int64{Value: 10}
	resp, err := cli.GetSession(context.Background(), sessionId)
	if err != nil {
		printl(err.Error())
		return
	}
	printj(resp)
}

func TestRpcGetSessions(t *testing.T) {
	userId := &rpc.Int64{Value: 1}
	stream, err := cli.GetSessions(context.Background(), userId)
	if err != nil {
		printl(err.Error())
		return
	}
	for {
		session, errr := stream.Recv()
		if errr != nil {
			if errr == io.EOF {
				break
			} else {
				printl(errr.Error())
				continue
			}
		}
		printj(session)
	}
}

func TestRpcGetMembers(t *testing.T) {
	sessionId := &rpc.Int64{Value: 10}
	stream, err := cli.GetMembers(context.Background(), sessionId)
	if err != nil {
		printl(err.Error())
		return
	}
	for {
		member, errr := stream.Recv()
		if errr != nil {
			if errr == io.EOF {
				break
			} else {
				printl(errr.Error())
				continue
			}
		}
		printj(member)
	}
}

func TestSendAndWithdraw(t *testing.T){
	TestRpcReceiveMessage(t)
	TestRpcWithdrawMessage(t)
}

func TestRpcReceiveMessage(t *testing.T) {
	words := rpc.Words{Text: "你好好大家好",}
	body, err := proto.Marshal(&words)
	if err != nil {
		printl(err.Error())
		return
	}
	message := &rpc.Message{
		SendId:    2,
		SessionId: 9,
		RequestId: 1116,
		Status:    rpc.MsgStatus_MS_NORMAL,
		Type:      rpc.MsgType_MT_WORDS,
		Content:   body,
	}
	_, err = cli.ReceiveMessage(context.Background(), message)
	if err != nil {
		printl(err.Error())
		return
	}
}

func TestRpcWithdrawMessage(t *testing.T) {
	req := &rpc.WithdrawMessageReq{
		SenderId:  2,
		MessageId:    33,
		SessionId: 9,
	}
	ret, err := cli.WithdrawMessage(context.Background(), req)
	if err != nil {
		printl(err.Error())
		return
	}
	printj(ret)
}

func TestRpcReceiveAck(t *testing.T) {
	ack := &rpc.Ack{
		ObjId: 20,
		Type:  rpc.AckType_AT_MESSAGE,
	}
	_, err := cli.ReceiveACK(context.Background(), ack)
	if err != nil {
		printl(err.Error())
		return
	}
}

func TestRpcSyncMessage(t *testing.T) {
	req := &rpc.SyncMessageReq{
		DeviceId:    18,
	}
	stream, err := cli.SyncMessages(context.Background(), req)
	if err != nil {
		printl(err.Error())
		return
	}
	for {
		message, errr := stream.Recv()
		if errr != nil {
			if errr == io.EOF {
				break
			} else {
				printl(errr.Error())
				continue
			}
		}
		printj(message)
	}
}

func TestAuthorization(t *testing.T) {
	req := &rpc.AuthReq{
		Code: "abc",
	}
	resp, err := cli.Authorization(context.Background(), req)
	if err!=nil{
		printl(err.Error())
		return
	}
	printj(resp)
}
