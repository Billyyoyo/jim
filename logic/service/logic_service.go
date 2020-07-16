package service

import (
	"context"
	"jim/common/rpc"
)

type LogicService struct {
}

// 连接建立后第一步，保存到设备库，将设备连接缓存到redis，更新设备表的最后连接时间
func (s *LogicService) Register(ctx context.Context, req *rpc.RegisterReq) (resp *rpc.RegisterResp, err error) {
	return
}

// 建立连接后第二步，获取所有会话，私聊和群组
func (s *LogicService) GetSessions(req *rpc.Int64, srv rpc.LogicService_GetSessionsServer) (err error) {
	return
}

// 创建群组
func (s *LogicService) CreateSession(ctx context.Context, req *rpc.CreateSessionReq) (session *rpc.Session, err error) {
	return
}

// 加入群组
func (s *LogicService) JoinSession(ctx context.Context, req *rpc.JoinSessionReq) (ret *rpc.Bool, err error) {
	return
}

// 退出群组
func (s *LogicService) QuitSession(ctx context.Context, req *rpc.QuitSessionReq) (ret *rpc.Bool, err error) {
	return
}

// 重命名群
func (s *LogicService) RenameSession(ctx context.Context, req *rpc.RenameSessionReq) (ret *rpc.Bool, err error) {
	return
}

// 当逻辑服务器收到来自连接服务器的用户消息传递
// 1.取得session下所有的member，并取得每个member的sequence
// 2.保存消息到数据库
// 3.发送客户端ack，并将ack保存到表中，记录发送数量
// 4.从缓存中取得所有用户的在线连接，直接传递消息给在线连接
// 5.不在线的用户无视 todo 思考
func (s *LogicService) HandleMessage(ctx context.Context, req *rpc.Message) (empty *rpc.Empty, err error) {
	return
}

// 发送给客户端的消息的ack处理  记录送达数量
func (s *LogicService) HandleACK(ctx context.Context, req *rpc.Ack) (empty *rpc.Empty, err error) {
	return
}

// 获取会话中的所有成员
func (s *LogicService) GetMember(req *rpc.Int64, srv rpc.LogicService_GetMemberServer) (err error) {
	return
}

// 同步消息
//1. 不传seq时按服务端自己保存的最后序列号往后查询
//2. seq为int64  获取单条
//3. seq有，用in获取
//4. seq有>，按范围获取
func (s *LogicService) SyncMessage(req *rpc.SyncMessageReq, srv rpc.LogicService_SyncMessageServer) (err error) {
	return
}

// 撤回
func (s *LogicService) WithdrawMessage(ctx context.Context, req *rpc.WithdrawMessageReq) (ret *rpc.Bool, err error) {
	return
}

// 离线
func (s *LogicService) Offline(ctx context.Context, req *rpc.OfflineReq) (empty *rpc.Empty, err error) {
	return
}
