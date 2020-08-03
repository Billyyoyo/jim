package handler

import (
	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
	"jim/common/rpc"
	"jim/common/tool"
	"jim/common/utils"
	"jim/logic/cache"
	"jim/logic/dao"
	"jim/logic/model"
)

func CreateSession(ownerId int64, name string, _type int8, userIds []int64) (session *model.Session, members *[]model.User, err error) {
	ss := dao.DB()
	defer ss.Close()
	err = ss.Begin()
	if err != nil {
		log.Error("create session - start transaction fail:", err.Error())
		return
	}
	owner, err := dao.GetUser(ownerId)
	if err != nil {
		log.Error("create session - get owner fail:", err.Error())
		ss.Rollback()
		return
	}
	// 创建会话
	session = &model.Session{
		Name:       name,
		Type:       _type,
		Owner:      ownerId,
		CreateTime: utils.GetCurrentMS(),
	}
	err = dao.CreateSession(ss, session)
	if err != nil {
		log.Error("create session - create session fail:", err.Error())
		ss.Rollback()
		return
	}
	members = &[]model.User{}
	// 创建者加入成员
	om := &model.Member{
		SessionId:  session.Id,
		UserId:     owner.Id,
		CreateTime: utils.GetCurrentMS(),
	}
	err = dao.AddMemberV2(ss, om)
	if err != nil {
		log.Error("create session - add owner fail:", err.Error())
		ss.Rollback()
		return
	}
	*members = append(*members, *owner)
	//加入其他成员
	for _, userId := range userIds {
		m := &model.Member{
			SessionId:  session.Id,
			UserId:     userId,
			CreateTime: utils.GetCurrentMS(),
		}
		user, errr := dao.GetUser(userId)
		if errr != nil {
			log.Error("create session - get user fail:", errr.Error())
			continue
		}
		errr = dao.AddMemberV2(ss, m)
		if errr != nil {
			log.Error("create session - add member fail:", errr.Error())
			continue
		}
		*members = append(*members, *user)
	}
	ss.Commit()
	// 通知成员会话创建成功
	tool.AsyncRun(func() {
		createSessionNext(session, members)
	})
	return
}

func createSessionNext(session *model.Session, users *[]model.User) {
	csa := &rpc.CreateSessionAction{
		SessionId: session.Id,
		Name:      session.Name,
		OwnerId:   session.Owner,
		Type:      rpc.SessionType(session.Type),
	}
	bs, err := proto.Marshal(csa)
	if err != nil {
		log.Error("create session - serial action fail:", err.Error())
		return
	}
	for _, user := range *users {
		// 检查用户是否在线
		conns, errr := cache.ListUserConn(user.Id)
		if errr != nil {
			log.Error("create session - load user conns fail:", errr.Error())
			continue
		}
		// 用户多端设备
		for _, conn := range *conns {
			action := &rpc.Action{
				UserId:     user.Id,
				RemoteAddr: conn.Addr,
				Time:       session.CreateTime,
				Type:       rpc.ActType_ACT_CREATE,
				Data:       bs,
			}
			ret := SendAction(conn.Server, action)
			if ret == 1 {
				cache.RemoveUserConn(user.Id, conn.DeviceId)
			}
		}
	}
}

func JoinSession(userId, sessionId int64) (err error) {
	user, err := dao.GetUser(userId)
	if err != nil {
		log.Error("join session - get user fail:", err.Error())
		return
	}
	member := &model.Member{
		SessionId:  sessionId,
		UserId:     user.Id,
		CreateTime: utils.GetCurrentMS(),
	}
	err = dao.AddMember(member)
	if err != nil {
		log.Error("join session - add member fail:", err.Error())
		return
	}
	tool.AsyncRun(func() {
		joinSessionNext(sessionId, user)
	})
	return
}

func joinSessionNext(sessionId int64, newbie *model.User) {
	newbier := rpc.User{
		Id:   newbie.Id,
		Name: newbie.Name,
	}
	jsa := rpc.JoinSessionAction{
		SessionId: sessionId,
		User:      &newbier,
	}
	bs, err := proto.Marshal(&jsa)
	if err != nil {
		log.Error("join session - serial action fail:", err.Error())
		return
	}
	users, err := dao.GetMemberInSession(sessionId)
	if err != nil {
		log.Error("join session - get members fail:", err.Error())
	}
	for _, user := range *users {
		// 检查用户是否在线
		conns, errr := cache.ListUserConn(user.Id)
		if errr != nil {
			log.Error("join session - load user conns fail:", errr.Error())
			continue
		}
		// 用户多端设备
		for _, conn := range *conns {
			action := &rpc.Action{
				UserId:     user.Id,
				RemoteAddr: conn.Addr,
				Time:       utils.GetCurrentMS(),
				Type:       rpc.ActType_ACT_JOIN,
				Data:       bs,
			}
			ret := SendAction(conn.Server, action)
			if ret == 1 {
				cache.RemoveUserConn(user.Id, conn.DeviceId)
			}
		}
	}
}

func QuitSession(userId, sessionId int64) (err error) {
	user, err := dao.GetUser(userId)
	if err != nil {
		log.Error("quit session - get user fail:", err.Error())
		return
	}
	err = dao.DeleteMember(sessionId, userId)
	if err != nil {
		log.Error("quit session - del member fail:", err.Error())
		return
	}
	tool.AsyncRun(func() {
		quitSessionNext(sessionId, user)
	})
	return
}

func quitSessionNext(sessionId int64, deler *model.User) {
	jsa := rpc.QuitSessionAction{
		SessionId: sessionId,
		UserId:    deler.Id,
	}
	bs, err := proto.Marshal(&jsa)
	if err != nil {
		log.Error("quit session - serial action fail:", err.Error())
		return
	}
	users, err := dao.GetMemberInSession(sessionId)
	if err != nil {
		log.Error("quit session - get members fail:", err.Error())
	}
	// 退出者也要发送
	*users = append(*users, *deler)
	for _, user := range *users {
		// 检查用户是否在线
		conns, errr := cache.ListUserConn(user.Id)
		if errr != nil {
			log.Error("quit session - load user conns fail:", errr.Error())
			continue
		}
		// 用户多端设备
		for _, conn := range *conns {
			action := &rpc.Action{
				UserId:     user.Id,
				RemoteAddr: conn.Addr,
				Time:       utils.GetCurrentMS(),
				Type:       rpc.ActType_ACT_QUIT,
				Data:       bs,
			}
			ret := SendAction(conn.Server, action)
			if ret == 1 {
				cache.RemoveUserConn(user.Id, conn.DeviceId)
			}
		}
	}
}

func RenameSession(userId, sessionId int64, name string) (err error) {
	isExist, err := dao.IsUserInSession(userId, sessionId)
	if err != nil {
		log.Error("rename session - check member fail:", err.Error())
		return
	}
	if !isExist {
		log.Error("rename session - not member")
		return
	}
	err = dao.RenameSession(sessionId, name)
	if err != nil {
		log.Error("rename session - rename fail:", err.Error())
		return
	}
	tool.AsyncRun(func() {
		renameSessionNext(sessionId, name)
	})
	return
}

func renameSessionNext(sessionId int64, name string) {
	jsa := rpc.RenameSessionAction{
		SessionId: sessionId,
		Name:      name,
	}
	bs, err := proto.Marshal(&jsa)
	if err != nil {
		log.Error("rename session - serial action fail:", err.Error())
		return
	}
	users, err := dao.GetMemberInSession(sessionId)
	if err != nil {
		log.Error("rename session - get members fail:", err.Error())
	}
	for _, user := range *users {
		// 检查用户是否在线
		conns, errr := cache.ListUserConn(user.Id)
		if errr != nil {
			log.Error("rename session - load user conns fail:", errr.Error())
			continue
		}
		// 用户多端设备
		for _, conn := range *conns {
			action := &rpc.Action{
				UserId:     user.Id,
				RemoteAddr: conn.Addr,
				Time:       utils.GetCurrentMS(),
				Type:       rpc.ActType_ACT_RENAME,
				Data:       bs,
			}
			ret := SendAction(conn.Server, action)
			if ret == 1 {
				cache.RemoveUserConn(user.Id, conn.DeviceId)
			}
		}
	}
}
