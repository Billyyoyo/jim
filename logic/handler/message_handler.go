package handler

import (
	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
	"jim/common/rpc"
	"jim/common/utils"
	"jim/logic/cache"
	"jim/logic/dao"
	"jim/logic/model"
	"strconv"
	"strings"
)

// 当逻辑服务器收到来自连接服务器的用户消息传递
func ReceiveMessage(sendId, sessionId, requestId int64, _type int8, content []byte) (err error) {
	// 获取cache中该会话的发送消息序列号
	sendNo, err := cache.GetSessionMsgSendNo(sessionId)
	if err != nil {
		log.Error("recieve message - get user send no fail: ", err.Error())
		return
	}
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
			SendNo:     sendNo,
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
			SendNo:     message.SendNo,
			SendId:     message.SenderId,
			SessionId:  message.SessionId,
			Time:       message.CreateTime,
			Status:     rpc.MsgStatus(message.Status),
			Type:       rpc.MsgType(message.Type),
			SequenceNo: message.Sequence,
			Content:    message.Body,
			DeviceId:   conn.DeviceId,
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
func SyncMessage(userId int64, cond string) (continuity bool, messages *[]model.Message, err error) {
	if no, err1 := strconv.ParseInt(cond, 10, 64); err1 == nil {
		messages, err = dao.GetMessagesSeqAfter(userId, no)
		if err != nil {
			log.Error("sync message - load after fail: ", err.Error())
			return
		}
		continuity = true
	} else if strings.Contains(cond, ",") {
		seqs := strings.Split(cond, ",")
		nos := []int64{}
		for _, seq := range seqs {
			no, _ := strconv.ParseInt(seq, 10, 64)
			nos = append(nos, no)
		}
		messages, err = dao.GetMessagesSeqIn(userId, nos)
		if err != nil {
			log.Error("sync message - load in fail: ", err.Error())
			return
		}
	} else if strings.Contains(cond, "<") {
		se := strings.Split(cond, "<")
		start, _ := strconv.ParseInt(se[0], 10, 64)
		end, _ := strconv.ParseInt(se[1], 10, 64)
		messages, err = dao.GetMessagesSeqRange(userId, start, end)
		if err != nil {
			log.Error("sync message - load range fail: ", err.Error())
			return
		}
		continuity = true
	}
	return
}

func WithdrawMessage(sessionId, userId, sendNo int64) (ret bool, err error) {
	// 先保存消息状态为撤回
	affect, err := dao.WithdrawMessage(userId, sendNo)
	if err != nil {
		log.Error("withdraw message - update message fail: ", err.Error())
		ret = false
		return
	}
	if affect == 0 {
		log.Error("withdraw message - no message updated")
		ret = false
		return
	}
	withdrawMessageNext(sessionId, userId, sendNo)
	return
}

func withdrawMessageNext(sessionId, userId, sendNo int64) {
	// 找出所有该消息的接收者
	receptorIds, err := dao.GetReceptorIdsInSendNo(sessionId, sendNo)
	if err != nil {
		log.Error("withdraw message - no receptor found: ", err.Error())
		return
	}
	for _, receptorId := range *receptorIds {
		wa := &rpc.WithdrawMessageAction{
			MessageId: 0,
			UserId:    userId,
		}
		bs, errr := proto.Marshal(wa)
		if errr != nil {
			log.Error("withdraw message - serial action content fail: ", err.Error())
			continue
		}
		conns, errr := cache.ListUserConn(receptorId)
		if errr != nil {
			log.Error("withdraw message - get user online connection fail: ", err.Error())
			continue
		}
		// 给所有接受者的device发送消息已撤回的动作
		for _, conn := range *conns {
			action := &rpc.Action{
				UserId:   receptorId,
				DeviceId: conn.DeviceId,
				Time:     utils.GetCurrentMS(),
				Type:     rpc.ActType_ACT_WITHDRAW,
				Data:     bs,
			}
			SendAction(conn.Server, action)
		}
	}
}
