package core

import (
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
	"jim/common/rpc"
	"jim/common/utils"
	"strconv"
	"strings"
)

func (cli *IMClient) sendMsg(content string) {
	if cli.Ctx.SessionId == 0 {
		fmt.Println("you don't enter any session")
		return
	}
	words := &rpc.Words{Text: content}
	body, err := proto.Marshal(words)
	if err != nil {
		log.Error("send msg - parse body err: ", err.Error())
		return
	}
	msg := &rpc.Message{
		SendId:    cli.Ctx.UserId,
		SessionId: cli.Ctx.SessionId,
		RequestId: utils.GetCurrentMS(),
		Type:      rpc.MsgType_MT_WORDS,
		Content:   body,
	}
	mbs, err := proto.Marshal(msg)
	if err != nil {
		log.Error("send msg - parse message err: ", err.Error())
		return
	}
	cli.sendPack(rpc.PackType_PT_MSG, &mbs)
}

func (cli *IMClient) CreateSession(content string) {
	arr := strings.Split(content, " ")
	if len(arr) < 3 {
		fmt.Println("params wrong, ex: name type uid1,uid2,uid3")
		return
	}
	ids := make([]int64, 0)
	idArr := strings.Split(arr[2], ",")
	for _, id := range idArr {
		if id64, err := strconv.ParseInt(id, 10, 64); err == nil {
			ids = append(ids, id64)
		}
	}
	if len(ids) == 0 {
		fmt.Println("member ids is wrong")
		return
	}
	_type, err := strconv.Atoi(arr[1])
	if err != nil {
		fmt.Println("session type is wrong. ex: person=1 group=2")
		return
	}
	act := &rpc.CreateSessionAction{
		OwnerId: cli.Ctx.UserId,
		Name:    arr[0],
		UserIds: ids,
		Type:    rpc.SessionType(_type),
	}
	mbs, err := proto.Marshal(act)
	if err != nil {
		log.Error("send create action - parse action err: ", err.Error())
		return
	}
	cli.sendAction(rpc.ActType_ACT_CREATE, &mbs)
}

func (cli *IMClient) JoinSession(content string) {
	arr := strings.Split(content, " ")
	if len(arr) < 1 {
		fmt.Println("params wrong, ex: :jo sid or :jo sid uid")
		return
	}
	sid, err := strconv.ParseInt(arr[0], 10, 64)
	if err != nil {
		fmt.Println("wrong sessionId")
		return
	}
	act := &rpc.JoinSessionAction{
		SessionId: sid,
		User:      &rpc.User{},
	}
	if len(arr) > 1 {
		act.User.Id, err = strconv.ParseInt(arr[1], 10, 64)
		if err != nil {
			fmt.Println("wrong userId")
			return
		}
	} else {
		act.User.Id = cli.Ctx.UserId
	}
	mbs, err := proto.Marshal(act)
	if err != nil {
		log.Error("send join action - parse action err: ", err.Error())
		return
	}
	cli.sendAction(rpc.ActType_ACT_JOIN, &mbs)
}

func (cli *IMClient) QuitSession(content string) {
	arr := strings.Split(content, " ")
	if len(arr) < 1 {
		fmt.Println("params wrong, ex: :qu sid or :jo sid uid")
		return
	}
	sid, err := strconv.ParseInt(arr[0], 10, 64)
	if err != nil {
		fmt.Println("wrong sessionId")
		return
	}
	act := &rpc.QuitSessionAction{
		SessionId: sid,
	}
	if len(arr) > 1 {
		act.UserId, err = strconv.ParseInt(arr[1], 10, 64)
		if err != nil {
			fmt.Println("wrong userId")
			return
		}
	} else {
		act.UserId = cli.Ctx.UserId
	}
	mbs, err := proto.Marshal(act)
	if err != nil {
		log.Error("send quit action - parse action err: ", err.Error())
		return
	}
	cli.sendAction(rpc.ActType_ACT_QUIT, &mbs)
}

func (cli *IMClient) RenameSession(content string) {
	arr := strings.Split(content, " ")
	if len(arr) < 2 {
		fmt.Println("params wrong, ex: :re sid new_name")
		return
	}
	sid, err := strconv.ParseInt(arr[0], 10, 64)
	if err != nil {
		fmt.Println("wrong sessionId")
		return
	}
	act := &rpc.RenameSessionAction{
		SessionId: sid,
		Name:      arr[1],
	}
	mbs, err := proto.Marshal(act)
	if err != nil {
		log.Error("send rename action - parse action err: ", err.Error())
		return
	}
	cli.sendAction(rpc.ActType_ACT_RENAME, &mbs)
}

func (cli *IMClient) WithdrawMsg(content string) {
	if cli.Ctx.SessionId == 0 {
		fmt.Println("you should enter a session")
		return
	}
	sendNo, err := strconv.ParseInt(content, 10, 64)
	if err != nil {
		fmt.Println("params wrong, ex: :- sendNo (sendNo is the index of message in session)")
		return
	}
	act := &rpc.WithdrawMessageAction{
		SessionId: cli.Ctx.SessionId,
		UserId:    cli.Ctx.UserId,
		SendNo:    sendNo,
	}
	mbs, err := proto.Marshal(act)
	if err != nil {
		log.Error("send withdraw action - parse action err: ", err.Error())
		return
	}
	cli.sendAction(rpc.ActType_ACT_WITHDRAW, &mbs)
}

func (cli *IMClient) SyncMsg() {
	act := &rpc.SyncMessageAction{
		UserId: cli.Ctx.UserId,
	}
	mbs, err := proto.Marshal(act)
	if err != nil {
		log.Error("sync msg - parse action err: ", err.Error())
		return
	}
	cli.sendAction(rpc.ActType_ACT_SYNC, &mbs)
}

func (cli *IMClient) GetSessions() {
	sessions, err := GetSessions(&cli.Ctx)
	if err != nil {
		log.Error("getsessions fail: ", err.Error())
		return
	}
	printj(sessions)
}

func (cli *IMClient) GetMembers() {
	if cli.Ctx.SessionId == 0 {
		fmt.Println("you should enter a session")
		return
	}
	members, err := GetMembers(&cli.Ctx)
	if err != nil {
		log.Error("getsessions fail: ", err.Error())
		return
	}
	printj(members)
}

func (cli *IMClient) SwitchSession(content string) {
	sessionId, err := strconv.ParseInt(content, 10, 64)
	if err != nil {
		fmt.Println("params wrong, ex: :sw sessionId")
		return
	}
	cli.Ctx.SessionId = sessionId
}

func printj(data interface{}) {
	bs, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		log.Error(err.Error())
		return
	}
	fmt.Println(string(bs))
}
