package handler

import (
	"errors"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"jim/common/rpc"
	"jim/common/utils"
	"jim/logic/cache"
	"jim/logic/dao"
	"jim/logic/model"
	"strconv"
	"strings"
)

func Authorization(code string) (uid int64, token string, err error) {
	// todo 需要到oauth2认证code有效
	if code == "" {
		err = errors.New("auth code is wrong")
		log.Error("authorization - auth code is wrong", err.Error())
		return
	}
	// todo 临时策略
	uid, err = strconv.ParseInt(code, 10, 64)
	if err != nil {
		log.Error("authorization - auth code is expired", err.Error())
		return
	}
	token = strings.ReplaceAll(uuid.New().String(), "-", "")
	err = cache.SaveUserToken(uid, token)
	return
}

func Validate(userId int64, deviceId int64, token string) (err error) {
	conn := &model.UserState{}
	err = cache.GetUserConn(userId, deviceId, conn)
	if err != nil {
		log.Error("validate - get user info fail:", err.Error())
		return
	}
	if token != conn.Token {
		log.Error("validate - compare token fail:", err.Error())
		return
	}
	return
}

func Register(userId int64, token, addr, server, serialNo string) (deviceId int64, err error) {
	// 检查token是否有效
	uid, err := cache.HasUserToken(token)
	if err != nil {
		log.Error("register - check token fail:", err.Error())
		return
	}
	if uid == 0 {
		log.Error("register - no this token:", err.Error())
		return
	}
	if uid != userId {
		log.Error("register - token is not belong this user:", err.Error())
		return
	}
	// 检查用户设备是否入库
	existInDB, device, err := dao.GetDevice(userId, serialNo)
	if err != nil {
		log.Error("register - get device in db fail:", err.Error())
		return
	}
	if !existInDB {
		device = &model.Device{
			UserId:       userId,
			SerialNo:     serialNo,
			LastConnTime: utils.GetCurrentMS(),
			LastAddress:  addr,
			CreateTime:   utils.GetCurrentMS(),
		}
	} else {
		device.LastAddress = addr
		device.LastConnTime = utils.GetCurrentMS()
		// 检查是否在线
		existInCache, errr := cache.HasUserConn(userId, device.Id)
		if errr != nil {
			log.Error("register - check conn exist fail:", errr.Error())
			err = errr
			return
		}
		if existInCache {
			// 在线 将在线的连接踢下线
			userConn := &model.UserState{}
			if er2 := cache.GetUserConn(userId, deviceId, userConn); er2 == nil {
				SendKickoff(userConn.Server, &rpc.Text{Value: userConn.Addr})
			}
		}
	}
	err = dao.SaveDevice(device)
	if err != nil {
		log.Error("register - save device in db fail:", err.Error())
		return
	}
	//更新redis中的device
	conn := &model.UserState{
		Server:   server,
		Addr:     addr,
		DeviceId: device.Id,
		Token:    token,
	}
	err = cache.SaveUserConn(device.UserId, conn)
	if err != nil {
		log.Error("register - save user conn in redis fail:", err.Error())
		return
	}
	deviceId = device.Id
	return
}

func Offline(userId, deviceId int64) (err error) {
	// 删除用户设备的在线状态
	cache.RemoveUserConn(userId, deviceId)
	// 更新数据库表中的用户设备数据
	device, err := dao.GetDeviceById(deviceId)
	if err != nil {
		log.Error("offline - get device in db fail: ", err.Error())
		return
	}
	// 更新最后消息序号
	device.LastDisconTime = utils.GetCurrentMS()
	dao.SaveDevice(device)
	return
}

func GetMembers(sessionId int64) (members *[]model.User, err error) {
	members, err = dao.GetMemberInSession(sessionId)
	if err != nil {
		log.Error("getmembers - get members fail: ", err.Error())
	}
	return
}

func GetSessions(userId int64) (sessions *[]model.Session, err error) {
	sessions, err = dao.GetSessionByUser(userId)
	if err != nil {
		log.Error("getsessions - get sessions fail: ", err.Error())
	}
	return
}

func GetSession(sessionId int64) (session *model.Session, members *[]model.User, err error) {
	session, err = dao.GetSession(sessionId)
	if err != nil {
		log.Error("getsession - get sessions fail: ", err.Error())
	}
	members, err = dao.GetMemberInSession(sessionId)
	if err != nil {
		log.Error("getsession - get members fail: ", err.Error())
	}
	return
}
