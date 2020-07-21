package service

import (
	"context"
	log "github.com/sirupsen/logrus"
	"jim/common/rpc"
	"jim/logic/handler"
	"jim/logic/model"
)

type LogicService struct {
}

// 连接建立后第一步，保存到设备库，将设备连接缓存到redis，更新设备表的最后连接时间
func (s *LogicService) Register(ctx context.Context, req *rpc.RegisterReq) (resp *rpc.RegisterResp, err error) {
	deviceId, lastSequence, err := handler.Register(req.UserId, req.Token, req.Addr, req.Server, req.SerialNo)
	if err != nil {
		return
	}
	resp = &rpc.RegisterResp{
		DeviceId:     deviceId,
		LastSequence: lastSequence,
	}
	return
}

// 建立连接后第二步，获取所有会话，私聊和群组
func (s *LogicService) GetSessions(req *rpc.Int64, stream rpc.LogicService_GetSessionsServer) (err error) {
	sessions, err := handler.GetSessions(req.Value)
	if err != nil {
		return
	}
	for _, ms := range *sessions {
		rs := &rpc.Session{
			Id:    ms.Id,
			Name:  ms.Name,
			Type:  rpc.SessionType(ms.Type),
			Owner: ms.Owner,
		}
		errr := stream.Send(rs)
		if errr != nil {
			log.Error("stream send session:", errr.Error())
		}
	}
	return
}

// 创建群组
func (s *LogicService) CreateSession(ctx context.Context, req *rpc.CreateSessionReq) (resp *rpc.SessionResp, err error) {
	ms, mm, err := handler.CreateSession(req.Creater, req.Name, int8(req.Type), req.UserIds)
	if err != nil {
		return
	}
	resp = &rpc.SessionResp{}
	resp.Session = &rpc.Session{
		Id:    ms.Id,
		Name:  ms.Name,
		Type:  rpc.SessionType(ms.Type),
		Owner: ms.Owner,
	}
	resp.Members = []*rpc.User{}
	for _, m := range *mm {
		mem := &rpc.User{
			Id:   m.Id,
			Name: m.Name,
		}
		resp.Members = append(resp.Members, mem)
	}
	return
}

// 获取群组及成员
func (s *LogicService) GetSession(ctx context.Context, sessionId *rpc.Int64) (resp *rpc.SessionResp, err error) {
	ms, mm, err := handler.GetSession(sessionId.Value)
	if err != nil {
		return
	}
	resp = &rpc.SessionResp{}
	resp.Session = &rpc.Session{
		Id:    ms.Id,
		Name:  ms.Name,
		Type:  rpc.SessionType(ms.Type),
		Owner: ms.Owner,
	}
	resp.Members = []*rpc.User{}
	for _, m := range *mm {
		mem := &rpc.User{
			Id:   m.Id,
			Name: m.Name,
		}
		resp.Members = append(resp.Members, mem)
	}
	return
}

// 加入群组
func (s *LogicService) JoinSession(ctx context.Context, req *rpc.JoinSessionReq) (ret *rpc.Bool, err error) {
	err = handler.JoinSession(req.UserId, req.SessionId)
	if err != nil {
		ret = &rpc.Bool{
			Value: false,
		}
	} else {
		ret = &rpc.Bool{
			Value: true,
		}
	}
	return
}

// 退出群组
func (s *LogicService) QuitSession(ctx context.Context, req *rpc.QuitSessionReq) (ret *rpc.Bool, err error) {
	err = handler.QuitSession(req.UserId, req.SessionId)
	if err != nil {
		ret = &rpc.Bool{
			Value: false,
		}
	} else {
		ret = &rpc.Bool{
			Value: true,
		}
	}
	return
}

// 重命名群
func (s *LogicService) RenameSession(ctx context.Context, req *rpc.RenameSessionReq) (ret *rpc.Bool, err error) {
	err = handler.RenameSession(req.UserId, req.SessionId, req.Name)
	if err != nil {
		ret = &rpc.Bool{
			Value: false,
		}
	} else {
		ret = &rpc.Bool{
			Value: true,
		}
	}
	return
}

// 接收消息
func (s *LogicService) ReceiveMessage(ctx context.Context, req *rpc.Message) (empty *rpc.Empty, err error) {
	_type := int8(req.Type)
	err = handler.ReceiveMessage(req.SendId, req.SessionId, req.RequestId, _type, req.Content)
	empty = &rpc.Empty{}
	return
}

// 发送给客户端的消息的ack处理  记录送达数量
func (s *LogicService) ReceiveACK(ctx context.Context, req *rpc.Ack) (empty *rpc.Empty, err error) {
	_type := int8(req.Type)
	err = handler.ReceiveAck(req.ObjId, _type)
	empty = &rpc.Empty{}
	return
}

// 获取会话中的所有成员
func (s *LogicService) GetMembers(req *rpc.Int64, stream rpc.LogicService_GetMembersServer) (err error) {
	members, err := handler.GetMembers(req.Value)
	for _, member := range *members {
		user := &rpc.User{
			Id:   member.Id,
			Name: member.Name,
		}
		errr := stream.Send(user)
		if errr != nil {
			log.Error("stream send member:", errr.Error())
		}
	}
	return
}

// 同步消息
func (s *LogicService) SyncMessage(req *rpc.SyncMessageReq, stream rpc.LogicService_SyncMessageServer) (err error) {
	continuity, msgs, err := handler.SyncMessage(req.UserId, req.Condition)
	if err != nil {
		return
	}
	var prevSeq int64
	for _, msg := range *msgs {
		// 检查消息是否连续有序
		if continuity && prevSeq > 0 && msg.Sequence-prevSeq > 1 {
			for i := prevSeq + 1; i < msg.Sequence; i++ {
				m := &rpc.Message{SequenceNo: i}
				errr := stream.Send(m)
				if errr != nil {
					log.Error("stream sync empty message:", errr.Error())
				}
			}
		}
		message := &rpc.Message{
			Id:         msg.Id,
			SendId:     msg.SenderId,
			SessionId:  msg.SessionId,
			Time:       msg.CreateTime,
			Status:     rpc.MsgStatus(msg.Status),
			Type:       rpc.MsgType(msg.Type),
			SequenceNo: msg.Sequence,
			SendNo:     msg.SendNo,
		}
		if msg.Status == model.MESSAGE_STATUS_NORMAL {
			message.Content = msg.Body
		} else {
			message.Content = nil
		}
		prevSeq = msg.Sequence
		errr := stream.Send(message)
		if errr != nil {
			log.Error("stream sync message:", errr.Error())
		}
	}
	return
}

// 撤回
func (s *LogicService) WithdrawMessage(ctx context.Context, req *rpc.WithdrawMessageReq) (ret *rpc.Bool, err error) {
	ok, err := handler.WithdrawMessage(req.SessionId, req.SenderId, req.SendNo)
	ret = &rpc.Bool{Value: ok}
	return
}

// 离线 删除redis中设备连接信息 更新数据库中设备的断开时间
func (s *LogicService) Offline(ctx context.Context, req *rpc.OfflineReq) (empty *rpc.Empty, err error) {
	err = handler.Offline(req.UserId, req.DeviceId, req.LastSequence)
	empty = &rpc.Empty{}
	return
}
