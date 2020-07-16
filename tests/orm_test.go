package tests

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"jim/common/utils"
	"jim/logic/dao"
	"jim/logic/model"
	"testing"
)

func TestGetUser(t *testing.T) {
	user, err := dao.GetUser(9)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	print(user)
}

func TestGetSessions(t *testing.T) {
	sessions, err := dao.GetSessionByUser(1)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	print(sessions)
}

func TestGetMembers(t *testing.T) {
	members, err := dao.GetMemberInSession(1)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	print(members)
}

func TestAddMember(t *testing.T) {
	member := &model.Member{
		SessionId:  1,
		UserId:     1,
		CreateTime: utils.GetCurrentMS(),
	}
	err := dao.AddMember(member)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func TestGetDevice(t *testing.T) {
	device, err := dao.GetDevice(1, "123adskf23ek2jrh")
	if err != nil {
		fmt.Println(err)
		return
	}
	print(device)
}

func TestRecordDevice(t *testing.T) {
	device := &model.Device{
		UserId:         1,
		SerialNo:       "123adskf23ek2jrh",
		LastAddress:    "10.8.240.133:42123",
		LastConnTime:   utils.GetCurrentMS(),
		CreateTime:     utils.GetCurrentMS(),
		LastDisconTime: 0,
	}
	err := dao.SaveDevice(device)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func TestCreateSession(t *testing.T) {
	session := &model.Session{
		Name:       "Go space",
		Type:       2,
		Owner:      1,
		CreateTime: utils.GetCurrentMS(),
	}
	dao.CreateSession(session)
	print(session)
}

func TestAddAck(t *testing.T) {
	ack := &model.Ack{
		MsgId:       1,
		SendCount:   1,
		ArriveCount: 0,
	}
	dao.AddAck(ack)
	print(ack)
}

func TestAddMessage(t *testing.T) {
	body := map[string]interface{}{
		"url":      "http://www.sina.com.cn",
		"duration": "10",
	}
	bs, err := json.Marshal(body)
	if err != nil {
		log.Error(err.Error())
		return
	}
	msg := &model.Message{
		SenderId:   1,
		SessionId:  1,
		Type:       2,
		Status:     1,
		DeviceId:   1,
		Sequence:    3,
		ReceptorId: 2,
		Body:       bs,
		CreateTime: utils.GetCurrentMS(),
	}
	dao.AddMessage(msg)
	print(msg)
}

func TestGetMessageList(t *testing.T) {
	msgs, err := dao.GetMessageByUserAndSequence(2, 1)
	if err != nil {
		log.Error(err.Error())
		return
	}
	for _, msg := range *msgs {
		body := map[string]interface{}{}
		err := json.Unmarshal(msg.Body, &body)
		if err != nil {
			log.Error(err.Error())
			continue
		}
		print(body)
	}
}

func TestDelMember(t *testing.T) {
	err := dao.DeleteMember(1, 3)
	if err != nil {
		log.Println(err.Error())
	}
}

func TestRenameSession(t *testing.T) {
	err := dao.RenameSession(1, "Happy Coder")
	if err != nil {
		log.Println(err.Error())
	}
}

func TestWithdrawMessage(t *testing.T){
	err:= dao.WithdrawMessage(1)
	if err != nil {
		log.Error(err.Error())
	}
}
