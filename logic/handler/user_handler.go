package handler

import "jim/logic/model"

func Register(userId, token, addr, server, serialNo string) (ret bool, deviceId, lastSequence int64, err error) {
	return
}

func Offline(userId, deviceId, lastSequence int64) {

}

func GetMembers(sessionId int64) (members *[]model.Member, err error) {
	return
}

func GetSessions(userId int64) (sessions *[]model.Session, err error) {
	return
}
