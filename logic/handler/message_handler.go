package handler

import (
	log "github.com/sirupsen/logrus"
	"jim/common/rpc"
	"jim/common/utils"
	"jim/logic/cache"
	"jim/logic/dao"
	"jim/logic/model"
)

// 当逻辑服务器收到来自连接服务器的用户消息传递
// 3.发送客户端ack，并将ack保存到表中，记录发送数量
// 4.从缓存中取得所有用户的在线连接，直接传递消息给在线连接
// 5.不在线的用户无视 todo 思考
func ReceiveMessage(sendId, sessionId, requestId int64, _type int8, content []byte) (err error) {
	// 查询所有session中的用户
	members, err := dao.GetMemberInSession(sessionId)
	if err != nil {
		log.Error("recieve message - load members fail: ", err.Error())
		return
	}
	for _, member := range *members {
		// 取得session下所有的member，并取得每个member的下一个sequence
		// todo 如果获取序列号成功 但是保存消息失败 那序列号将无法回退
		sequence, errr := cache.GetUserMsgSequence(member.Id)
		if errr != nil {
			log.Error("recieve message - get next sequence fail: ", errr.Error())
			continue
		}
		// 保存消息到数据库，落地
		message := &model.Message{
			SenderId:   sendId,
			SessionId:  sessionId,
			Type:       _type,
			Status:     model.MESSAGE_STATUS_NORMAL,
			Sequence:   sequence,
			ReceptorId: member.Id,
			Body:       content,
			CreateTime: utils.GetCurrentMS(),
		}
		errr = dao.AddMessage(message)
		if errr != nil {
			log.Error("recieve message - save message fail: ", errr.Error())
			continue
		}
		// 发送到接入服务器，发送到用户客户端
		receiveMessageNext(&member, message)
	}
	// 发送客户端ack  该流程交给接入服务器  rpc返回无error的时候发送服务器收到的ack给客户端
	return
}

// 可以go程来处理后续
func receiveMessageNext(member *model.User, message *model.Message) {
	ack := &model.Ack{MsgId: message.Id}
	conns, err := cache.ListUserConn(member.Id)
	if err != nil {
		log.Error("receive message - load all conns fail: ", err.Error())
		return
	}
	for _, conn := range *conns {
		ack.SendCount++
		msg := &rpc.Message{
			Id:         message.Id,
			SendId:     message.SenderId,
			SessionId:  message.SessionId,
			Time:       message.CreateTime,
			Status:     rpc.MsgStatus(message.Status),
			Type:       rpc.MsgType(message.Type),
			SequenceNo: message.Sequence,
			Content:    message.Body,
		}
		SendMessage(conn.Server, msg)
	}
	dao.AddAck(ack)
}

// 给客户端发送消息或推送通知之后，客户端接收到消息或通知将发送该消息的ack给服务端，服务端接收到ack后更新ack表统计消息到达量
func ReceiveAck(ackId int64, _type int8) (err error) {
	if _type == model.ACK_TYPE_MESSAGE {
		err = dao.AccumulateAckArriveCount(ackId)
	} else if _type == model.ACK_TYPE_NOTIFICATION {

	}
	return
}

// 当客户端发现消息序列号不连续的情况下，请求服务端同步消息
// 条件1.单条消息序列号
// 条件2.多条消息序列号，用,号连接
// 条件3.消息序列号范围，startNo<endNo
// 当服务端发现消息列表序列号不连续或不存在，以空消息填充（type=0）
func SyncMessage(userId int64, conditions []string) (messages *[]model.Message, err error) {
	
	return
}

func WithdrawMessage(userId, messageId int64) (err error) {

	return
}
