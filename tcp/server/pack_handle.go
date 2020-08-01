package server

import (
	"github.com/golang/protobuf/proto"
	"github.com/panjf2000/gnet"
	log "github.com/sirupsen/logrus"
	"jim/common/rpc"
	"jim/common/tool"
	"jim/common/utils"
	"time"
)

func (cs *TcpServer) handlePing(c gnet.Conn) {
	pingPack := &rpc.Output{Type: rpc.PackType_PT_PONG,}
	bs, err := proto.Marshal(pingPack)
	if err != nil {
		return
	}
	c.AsyncWrite(bs)
}

func (cs *TcpServer) handlePong(c gnet.Conn) {
	data := c.Context().(ConnData)
	data.PongTime = time.Now().Unix()
	c.SetContext(data)
	cs.uconns.Store(c.RemoteAddr().String(), c)
}

func (cs *TcpServer) handleReg(c gnet.Conn, pack *rpc.RegInfo) (err error) {
	// 当连接建立 调用logic的register
	// 同步注册成功 将deviceId放到c.SetContext()中
	// 验证用户连接信息
	did, ls, err := tool.RegisterConnection(pack.UserId, c.RemoteAddr().String(), c.LocalAddr().String(), pack.SerialNo, pack.Token)
	if err != nil {
		// 注册失败 向客户端发送失败指令
		pingPack := &rpc.Output{
			Type: rpc.PackType_PT_AUTH,
			Code: 1,
			Info: err.Error(),
		}
		var bs []byte
		bs, err = proto.Marshal(pingPack)
		if err != nil {
			return
		}
		c.AsyncWrite(bs)
		return
	}
	// 保存用户连接数据
	ctx := ConnData{
		Did:      did,
		Uid:      pack.UserId,
		Serial:   pack.SerialNo,
		Seq:      ls,
		PongTime: time.Now().Unix(),
	}
	c.SetContext(ctx)
	cs.uconns.Store(c.RemoteAddr().String(), c)
	// 向客户端发送完成注册指令
	didr := &rpc.Int64{Value: did}
	body, err := proto.Marshal(didr)
	if err != nil {
		return
	}
	pingPack := &rpc.Output{
		Type: rpc.PackType_PT_AUTH,
		Code: 0,
		Data: body,
	}
	bs, err := proto.Marshal(pingPack)
	if err != nil {
		return
	}
	c.AsyncWrite(bs)
	return
}

func (cs *TcpServer) handleMsg(c *gnet.Conn, msg *rpc.Message) {
	// 提交给逻辑服务器
	err := tool.SendMsg(msg)
	// 逻辑服务器消息落地之后  返回客户端一个收到回执
	ack := &rpc.Ack{
		Type:      rpc.AckType_AT_MESSAGE,
		RequestId: msg.RequestId,
	}
	body, errr := proto.Marshal(ack)
	if errr != nil {
		return
	}
	pack := &rpc.Output{
		Type: rpc.PackType_PT_ACK,
		Data: body,
	}
	if err != nil {
		pack.Code = 1
		pack.Info = err.Error()
	} else {
		pack.Code = 0
	}
	bs, errr := proto.Marshal(pack)
	if errr != nil {
		return
	}
	(*c).AsyncWrite(bs)
}

func (cs *TcpServer) handleAct(c *gnet.Conn, pack *rpc.Action) {
	if pack.Type == rpc.ActType_ACT_JOIN {
		act := &rpc.JoinSessionAction{}
		err := proto.Unmarshal(pack.Data, act)
		if err != nil {
			return
		}
		cs.handleActJoin(c, pack.RequestId, act)
	} else if pack.Type == rpc.ActType_ACT_QUIT {
		act := &rpc.QuitSessionAction{}
		err := proto.Unmarshal(pack.Data, act)
		if err != nil {
			return
		}
		cs.handleActQuit(c, pack.RequestId, act)
	} else if pack.Type == rpc.ActType_ACT_WITHDRAW {
		act := &rpc.WithdrawMessageAction{}
		err := proto.Unmarshal(pack.Data, act)
		if err != nil {
			return
		}
		cs.handleActWithdraw(c, pack.RequestId, act)
	} else if pack.Type == rpc.ActType_ACT_SYNC {
		act := &rpc.SyncMessageAction{}
		err := proto.Unmarshal(pack.Data, act)
		if err != nil {
			return
		}
		cs.handleActSync(c, pack.RequestId, act)
	} else if pack.Type == rpc.ActType_ACT_RENAME {
		act := &rpc.RenameSessionAction{}
		err := proto.Unmarshal(pack.Data, act)
		if err != nil {
			return
		}
		cs.handleActRename(c, pack.UserId, pack.RequestId, act)
	} else if pack.Type == rpc.ActType_ACT_CREATE {
		act := &rpc.CreateSessionAction{}
		err := proto.Unmarshal(pack.Data, act)
		if err != nil {
			return
		}
		cs.handleActCreate(c, pack.RequestId, act)
	}
}

func (cs *TcpServer) handleAck(c *gnet.Conn, ack *rpc.Ack) {
	err := tool.SendAck(ack)
	if err != nil {
		log.Info("send ack error: ", err.Error())
		return
	}
	// 收到客户端发来的ack 根据最后序列号大小   更新最后序列号
	if ack.Type == rpc.AckType_AT_MESSAGE {
		data := (*c).Context().(ConnData)
		if data.Seq < ack.Seq {
			data.Seq = ack.Seq
			(*c).SetContext(data)
		}
	}
}

func (cs *TcpServer) handleOffline(data *ConnData) {
	tool.Offline(data.Uid, data.Did, data.Seq)
}

func (cs *TcpServer) handleActJoin(c *gnet.Conn, requestId int64, act *rpc.JoinSessionAction) {
	ok, err := tool.JoinSession(act.User.Id, act.SessionId)
	if err != nil || !ok {
		cs.sendAck(c, requestId, 1, "failed")
	} else {
		cs.sendAck(c, requestId, 0, "")
	}
}

func (cs *TcpServer) handleActQuit(c *gnet.Conn, requestId int64, act *rpc.QuitSessionAction) {
	ok, err := tool.QuitSession(act.UserId, act.UserId)
	if err != nil || !ok {
		cs.sendAck(c, requestId, 1, "failed")
	} else {
		cs.sendAck(c, requestId, 0, "")
	}
}

func (cs *TcpServer) handleActWithdraw(c *gnet.Conn, requestId int64, act *rpc.WithdrawMessageAction) {
	ok, err := tool.WithdrawMsg(act.UserId, act.SessionId, act.SendNo)
	if err != nil || !ok {
		cs.sendAck(c, requestId, 1, "failed")
	} else {
		cs.sendAck(c, requestId, 0, "")
	}
}

func (cs *TcpServer) handleActCreate(c *gnet.Conn, requestId int64, act *rpc.CreateSessionAction) {
	_, err := tool.CreateSession(act.OwnerId, int8(act.Type), act.Name, act.UserIds)
	if err != nil {
		cs.sendAck(c, requestId, 1, "failed")
	} else {
		cs.sendAck(c, requestId, 0, "")
	}
}

func (cs *TcpServer) handleActSync(c *gnet.Conn, requestId int64, act *rpc.SyncMessageAction) {
	msgs, err := tool.SyncMsgs(act.UserId, act.Cond)
	if err != nil {
		cs.sendAck(c, requestId, 1, "failed")
		return
	}
	act.Messages = *msgs
	ba, err := proto.Marshal(act)
	if err != nil {
		log.Error("sync msg - marshal body fail: ", err.Error())
		return
	}
	action := &rpc.Action{
		UserId:    act.UserId,
		RequestId: requestId,
		Time:      utils.GetCurrentMS(),
		Type:      rpc.ActType_ACT_SYNC,
		Data:      ba,
	}
	bs, err := proto.Marshal(action)
	if err != nil {
		log.Error("sync msg - marshal act fail: ", err.Error())
		return
	}
	outPack := &rpc.Output{
		Type: rpc.PackType_PT_ACTION,
		Data: bs,
	}
	bb, err := proto.Marshal(outPack)
	if err != nil {
		log.Error("sync msg - marshal pack fail: ", err.Error())
		return
	}
	(*c).AsyncWrite(bb)
}

func (cs *TcpServer) handleActRename(c *gnet.Conn, userId, requestId int64, act *rpc.RenameSessionAction) {
	ok, err := tool.RenameSession(userId, act.SessionId, act.Name)
	if err != nil || !ok {
		cs.sendAck(c, requestId, 1, "failed")
	} else {
		cs.sendAck(c, requestId, 0, "")
	}
}

func (cs *TcpServer) sendAck(c *gnet.Conn, requestId int64, code int32, info string) {
	ack := &rpc.Ack{
		Type:      rpc.AckType_AT_ACT,
		RequestId: requestId,
	}
	body, errr := proto.Marshal(ack)
	if errr != nil {
		return
	}
	pack := &rpc.Output{
		Type: rpc.PackType_PT_ACK,
		Data: body,
	}
	pack.Code = code
	pack.Info = info

	bs, errr := proto.Marshal(pack)
	if errr != nil {
		return
	}
	(*c).AsyncWrite(bs)
}
