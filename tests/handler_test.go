package tests

import (
	"github.com/golang/protobuf/proto"
	"jim/common/rpc"
	"jim/logic/handler"
	"jim/logic/model"
	"math/rand"
	"testing"
)

func TestCreateSessionHandler(t *testing.T) {
	session, members, err := handler.CreateSession(1, "Test Create Session", model.SESSION_TYPE_GROUP, []int64{1, 2, 3})
	if err != nil {
		println(err.Error())
		return
	}
	print(session)
	print(members)
}

func TestJoinSessionHandler(t *testing.T) {
	err := handler.JoinSession(4, 9)
	if err != nil {
		println(err.Error())
		return
	}
}

func TestQuitSessionHandler(t *testing.T) {
	err := handler.QuitSession(4, 9)
	if err != nil {
		println(err.Error())
		return
	}
}

func TestRenameSessionHandler(t *testing.T) {
	err := handler.RenameSession(3, 9, "Hello World")
	if err != nil {
		println(err.Error())
		return
	}
}

func TestRegisterHandle(t *testing.T) {
	serialNo := "20139fcd-25fe-42e0-9457-49356018beb9" //uuid.New().String()
	deviceId, sequence, err := handler.Register(2, "123123123", "localhost:42402", "localhost:5000", serialNo)
	if err != nil {
		println(err.Error())
		return
	}
	println("register success: ", deviceId, "  /  ", sequence)
}

func TestOfflineHandle(t *testing.T) {
	err := handler.Offline(1, 16, 3)
	if err != nil {
		println(err.Error())
	}
}

func TestGetSessionsHandle(t *testing.T) {
	sessions, err := handler.GetSessions(1)
	if err != nil {
		println(err.Error())
		return
	}
	print(sessions)
}

func TestGetSessionHandle(t *testing.T) {
	session, members, err := handler.GetSession(9)
	if err != nil {
		println(err.Error())
		return
	}
	print(session)
	print(members)
}

func TestGetMembersHandle(t *testing.T) {
	members, err := handler.GetMembers(9)
	if err != nil {
		println(err.Error())
		return
	}
	print(members)
}

func TestReceiveMessageHandle(t *testing.T) {
	message := rpc.Words{Text: "Good morning",}
	body, err := proto.Marshal(&message)
	if err != nil {
		println(err.Error())
		return
	}
	requestId := rand.Int63n(1000000)
	err = handler.ReceiveMessage(2, 9, requestId, model.MESSAGE_TYPE_WORDS, body)
	if err != nil {
		println(err.Error())
	}
}

func TestWithDrawMessageHandle(t *testing.T) {
	ok, err := handler.WithdrawMessage(9, 2, 4)
	if err != nil {
		println(err.Error())
		return
	}
	println(ok)
}

func TestSyncMessageHandle(t *testing.T) {
	continuity, msgs, err := handler.SyncMessage(1, "1,4")
	if err != nil {
		println(err.Error())
		return
	}
	if continuity {
		println("需要检查连续性")
	}
	print(msgs)
}

func TestReceiveAck(t *testing.T) {
	err := handler.ReceiveAck(7, model.ACK_TYPE_MESSAGE)
	if err != nil {
		println(err.Error())
	}
}
