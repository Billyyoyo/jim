package rpc

import (
	"context"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"jim/common/rpc"
)

var (
	cli rpc.LogicServiceClient
)

// 暂时使用固定地址调用logic服务  todo 后面应该改为负载均衡
func init() {
	conn, err := grpc.Dial("localhost:4000", grpc.WithInsecure())
	if err != nil {
		panic("grpc start up error: " + err.Error())
		return
	}
	cli = rpc.NewLogicServiceClient(conn)
}

func CreateSession(creater int64, _type int8, name string, userIds []int64) (session *rpc.SessionResp, err error){
	var Type rpc.SessionType
	if _type==1{
		Type = rpc.SessionType_SESSION_PERSON
	} else if _type==2 {
		Type = rpc.SessionType_SESSION_GROUP
	} else{

	}
	req := &rpc.CreateSessionReq{
		Name:    name,
		Creater: creater,
		Type:    Type,
		UserIds: userIds,
	}
	session, err = cli.CreateSession(context.Background(), req)
	if err!=nil{
		log.Error("create session - ", err.Error())
		return
	}
	return
}