package service

import (
	"context"
	log "github.com/sirupsen/logrus"
	"jim/common/rpc"
	"jim/logic/handler"
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
func (s *LogicService) HandleMessage(ctx context.Context, req *rpc.Message) (empty *rpc.Empty, err error) {
	_type := int8(req.Type)
	err = handler.ReceiveMessage(req.SendId, req.SessionId, req.RequestId, _type, req.Content)
	empty = &rpc.Empty{}
	return
}

// 发送给客户端的消息的ack处理  记录送达数量
func (s *LogicService) HandleACK(ctx context.Context, req *rpc.Ack) (empty *rpc.Empty, err error) {
	_type := int8(req.Type)
	err = handler.ReceiveAck(req.ObjId, _type)
	empty = &rpc.Empty{}
	return
}

// 获取会话中的所有成员
func (s *LogicService) GetMember(req *rpc.Int64, stream rpc.LogicService_GetMemberServer) (err error) {
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
//1. 不传seq时按服务端自己保存的最后序列号往后查询
//2. seq为int64  获取单条
//3. seq有，用in获取
//4. seq有>，按范围获取
func (s *LogicService) SyncMessage(req *rpc.SyncMessageReq, stream rpc.LogicService_SyncMessageServer) (err error) {
	msgs, err := handler.SyncMessage(req.UserId, req.SeqRanges)
	if err != nil {
		return
	}
	for _, msg := range *msgs {
		message := &rpc.Message{
			Id:        msg.Id,
			SendId:    msg.SenderId,
			SessionId: msg.SessionId,
			Time:      msg.CreateTime,
			Status:    rpc.MsgStatus(msg.Status),
			Type:      rpc.MsgType(msg.Type),
			Content:   msg.Body,
		}
		errr := stream.Send(message)
		if errr != nil {
			log.Error("stream sync message:", errr.Error())
		}
	}
	return
}

// 撤回
func (s *LogicService) WithdrawMessage(ctx context.Context, req *rpc.WithdrawMessageReq) (ret *rpc.Bool, err error) {
	err = handler.WithdrawMessage(req.UserId, req.MessageId)
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

// 离线 删除redis中设备连接信息 更新数据库中设备的断开时间
func (s *LogicService) Offline(ctx context.Context, req *rpc.OfflineReq) (empty *rpc.Empty, err error) {
	err = handler.Offline(req.UserId, req.DeviceId, req.LastSequence)
	empty = &rpc.Empty{}
	return
}
