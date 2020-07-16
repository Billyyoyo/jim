package handler

import "jim/logic/model"

func ReceiveMessage(sendId, sessionId, deviceId, requestId int64, _type int8, content []byte) (err error) {
	return
}

func ReceiveAck(ackId int64, _type int8)(err error) {
	return
}

func SyncMessage(userId int64) (messages *[]model.Message, err error) {

	return
}

func WithdrawMessage(userId, messageId int64) (ret bool, err error) {

	return
}
