package service

import (
	"context"
	"jim/common/rpc"
)

type TransService struct {
}

func (s *TransService) SendMessage(ctx context.Context, req *rpc.Message) (empty *rpc.Empty, err error) {
	return
}

func (s *TransService) SendNotification(ctx context.Context, req *rpc.Notification) (empty *rpc.Empty, err error) {
	return
}

func (s *TransService) SendAction(ctx context.Context, req *rpc.Action) (empty *rpc.Empty, err error) {
	return
}
