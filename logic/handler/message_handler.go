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
	"strconv"
	"strings"
)

// 当逻辑服务器收到来自连接服务器的用户消息传递
func ReceiveMessage(sendId, sessionId, requestId int64, _type int8, content []byte) (msgId int64, err error) {
	// 获取cache中该会话的发送消息序列号
	sequence, err := cache.GetSessionMsgSequence(sessionId)
	if err != nil {
		log.Error("recieve message - get session sequence fail: ", err.Error())
		return
	}
	now := utils.GetCurrentMS()
	// 保存消息到存储表，落地
	message := model.Message{
		SenderId:   sendId,
		SessionId:  sessionId,
		Type:       _type,
		Status:     model.MESSAGE_STATUS_NORMAL,
		Sequence:   sequence,
		Body:       content,
		CreateTime: now,
	}
	err = dao.AddMessage(message)
	if err != nil {
		log.Error("recieve message - save message fail: ", err.Error())
		return
	}
	msgId = message.Id
	dao.AddAck(model.Ack{MsgId: message.Id})

	tool.AsyncRun(func() {
		receiveMessageNext(message)
	})
	// 发送客户端ack  该流程交给接入服务器  rpc返回无error的时候发送服务器收到的ack给客户端
	return
}

// go程来处理后续
func receiveMessageNext(msg model.Message) {
	devices, err := dao.GetDevicesInSession(msg.SessionId)
	if err != nil {
		log.Error("receive message - load all device id fail:", err.Error())
		return
	}
	now := utils.GetCurrentMS()
	for _, device := range devices {
		// 每个设备的离线消息落地
		omsg := model.OffLineMessage{
			DeviceId:   device.Id,
			MessageId:  msg.Id,
			CreateTime: now,
		}
		errr := dao.AddOfflineMessage(omsg)
		if errr != nil {
			log.Error("receive message - save offline msg fail: ", errr.Error())
			continue
		}
		// 检查设备在线状态
		state, errr := cache.GetUserConn(device.UserId, device.Id)
		if errr != nil {
			log.Error("receive message - load user conn fail: ", errr.Error())
			// 如果不在线  只保存离线消息  直接跳过发送
			continue
		}

		rmsg := rpc.Message{
			Id:         msg.Id,
			SendId:     msg.SenderId,
			SessionId:  msg.SessionId,
			Time:       msg.CreateTime,
			Status:     rpc.MsgStatus(msg.Status),
			Type:       rpc.MsgType(msg.Type),
			SequenceNo: msg.Sequence,
			Content:    msg.Body,
			RemoteAddr: state.Addr,
			RequestId:  omsg.Id,
		}
		ret := SendMessage(state.Server, rmsg)
		if ret == 1 {
			cache.RemoveUserConn(device.UserId, device.Id)
		}
	}
}

// 给客户端发送消息或推送通知之后，客户端接收到消息或通知将发送该消息的ack给服务端，服务端接收到ack后更新ack表统计消息到达量
func ReceiveAck(objId, reqId int64, _type int8) (err error) {
	if _type == model.ACK_TYPE_MESSAGE {
		_ = dao.AccumulateAckArriveCount(objId)
		// 删除落地的离线消息
		err = dao.DeleteOfflineMsg(reqId)
		if err != nil {
			log.Error("receive ack - delete offline msg fail: ", err.Error())
			return
		}
	} else if _type == model.ACK_TYPE_NOTIFICATION {

	}
	return
}

// 当客户端发现消息序列号不连续的情况下，请求服务端同步消息
// 条件1.单条消息序列号
// 条件2.多条消息序列号，用,号连接
// 条件3.消息序列号范围，startNo<endNo
// 当服务端发现消息列表序列号不连续或不存在，以空消息填充（type=0）
func SyncMessage(deviceId int64) (messages []model.OMessage, err error) {
	// 直接返回该设备的离线消息列表
	messages, err = dao.GetOfflineMsgs(deviceId)
	return
}

func ListSessionMessages(sessionId int64, cond string) (continuity bool, messages []model.Message, err error) {
	if strings.Contains(cond, ":") {
		se := strings.Split(cond, ":")
		start, _ := strconv.ParseInt(se[0], 10, 64)
		count, _ := strconv.Atoi(se[1])
		messages, err = dao.GetMessagesSeqAfter(sessionId, start, count)
		if err != nil {
			log.Error("list message - load after fail: ", err.Error())
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
		messages, err = dao.GetMessagesSeqIn(sessionId, nos)
		if err != nil {
			log.Error("list message - load in fail: ", err.Error())
			return
		}
	} else if strings.Contains(cond, "<") {
		se := strings.Split(cond, "<")
		start, _ := strconv.ParseInt(se[0], 10, 64)
		end, _ := strconv.ParseInt(se[1], 10, 64)
		messages, err = dao.GetMessagesSeqRange(sessionId, start, end)
		if err != nil {
			log.Error("list message - load range fail: ", err.Error())
			return
		}
		continuity = true
	} else {
		messages = []model.Message{}
	}
	return
}

func WithdrawMessage(sessionId, senderId, messageId int64) (ret bool, err error) {
	// 先保存消息状态为撤回
	affect, err := dao.WithdrawMessage(sessionId, senderId, messageId)
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
	ret = true
	tool.AsyncRun(func() {
		withdrawMessageNext(sessionId, senderId, messageId)
	})
	return
}

func withdrawMessageNext(sessionId, senderId, messageId int64) {
	// 找出所有该消息的接收者
	members, err := dao.GetMemberInSession(sessionId)
	if err != nil {
		log.Error("withdraw message - no receptor found: ", err.Error())
		return
	}
	for _, member := range members {
		wa := rpc.WithdrawMessageAction{
			MessageId: messageId,
			SessionId: sessionId,
			UserId:    senderId,
		}
		bs, errr := proto.Marshal(&wa)
		if errr != nil {
			log.Error("withdraw message - serial action content fail: ", err.Error())
			continue
		}
		conns, errr := cache.ListUserConn(member.Id)
		if errr != nil {
			log.Error("withdraw message - get user online connection fail: ", err.Error())
			continue
		}
		// 给所有接受者的device发送消息已撤回的动作
		for _, conn := range conns {
			action := rpc.Action{
				UserId:     member.Id,
				RemoteAddr: conn.Addr,
				Time:       utils.GetCurrentMS(),
				Type:       rpc.ActType_ACT_WITHDRAW,
				Data:       bs,
			}
			ret := SendAction(conn.Server, action)
			if ret == 1 {
				cache.RemoveUserConn(member.Id, conn.DeviceId)
			}
		}
	}
}
