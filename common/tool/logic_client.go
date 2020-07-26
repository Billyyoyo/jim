package tool

import (
	"context"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"io"
	"jim/common/rpc"
	"strconv"
)

var (
	cli rpc.LogicServiceClient
)

// 暂时使用固定地址调用logic服务  todo 后期应该改为负载均衡
func init() {
	conn, err := grpc.Dial("localhost:4000", grpc.WithInsecure())
	if err != nil {
		panic("grpc start up error: " + err.Error())
		return
	}
	cli = rpc.NewLogicServiceClient(conn)
}

func Authorization(userId int64, code string) (ret bool, token string) {
	req := &rpc.AuthReq{
		Uid:  userId,
		Code: code,
	}
	resp, err := cli.Authorization(context.Background(), req)
	if err != nil || !resp.Ret {
		ret = false
	} else {
		ret = true
		token = resp.Token
	}
	return
}

func Validate(userId, deviceId, token string) bool {
	uid, err := strconv.ParseInt(userId, 10, 64)
	if err != nil {
		return false
	}
	did, err := strconv.ParseInt(deviceId, 10, 64)
	if err != nil {
		return false
	}
	if uid == 0 || did == 0 || token == "" {
		return false
	}
	req := &rpc.ValidReq{
		Uid:      uid,
		DeviceId: did,
		Token:    token,
	}
	resp, err := cli.Validate(context.Background(), req)
	if err != nil || !resp.Ret {
		return false
	} else {
		return true
	}
}

func CreateSession(creater int64, _type int8, name string, userIds []int64) (session *rpc.SessionResp, err error) {
	var Type rpc.SessionType
	if _type == 1 {
		Type = rpc.SessionType_SESSION_PERSON
	} else if _type == 2 {
		Type = rpc.SessionType_SESSION_GROUP
	} else {

	}
	req := &rpc.CreateSessionReq{
		Name:    name,
		Creater: creater,
		Type:    Type,
		UserIds: userIds,
	}
	session, err = cli.CreateSession(context.Background(), req)
	if err != nil {
		log.Error("create session - ", err.Error())
		return
	}
	return
}

func GetSessions(userId int64) (sessions *[]rpc.Session, err error) {
	stream, err := cli.GetSessions(context.Background(), &rpc.Int64{Value: userId})
	if err != nil {
		log.Error("GetSessions - rpc call:", err.Error())
		return
	}
	sessions = &[]rpc.Session{}
	for {
		session, errr := stream.Recv()
		if errr != nil {
			if errr == io.EOF {
				break
			} else {
				log.Error("GetSessions - rpc receive:", errr.Error())
				continue
			}
		}
		*sessions = append(*sessions, *session)
	}
	return
}

func GetMembers(sessionId int64) (members *[]rpc.User, err error) {
	stream, err := cli.GetMembers(context.Background(), &rpc.Int64{Value: sessionId})
	if err != nil {
		log.Error("GetMembers - rpc call:", err.Error())
		return
	}
	members = &[]rpc.User{}
	for {
		member, errr := stream.Recv()
		if errr != nil {
			if errr == io.EOF {
				break
			} else {
				log.Error("GetMembers - rpc receive:", errr.Error())
				continue
			}
		}
		*members = append(*members, *member)
	}
	return
}

func GetSession(sessionId int64) (resp *rpc.SessionResp, err error) {
	resp, err = cli.GetSession(context.Background(), &rpc.Int64{Value: sessionId})
	if err != nil {
		log.Error("GetSession - rpc call:", err.Error())
	}
	return
}

func JoinSession(userId, sessionId int64) (ret bool, err error) {
	req := &rpc.JoinSessionReq{
		UserId:    userId,
		SessionId: sessionId,
	}

	retR, err := cli.JoinSession(context.Background(), req)
	if err != nil {
		log.Error("JoinSession - rpc call:", err.Error())
		ret = false
		return
	}
	ret = retR.Value
	return
}

func QuitSession(userId, sessionId int64) (ret bool, err error) {
	req := &rpc.QuitSessionReq{
		UserId:    userId,
		SessionId: sessionId,
	}
	retR, err := cli.QuitSession(context.Background(), req)
	if err != nil {
		log.Error("QuitSession - rpc call:", err.Error())
		ret = false
		return
	}
	ret = retR.Value
	return
}

func RenameSession(userId, sessionId int64, name string) (ret bool, err error) {
	req := &rpc.RenameSessionReq{
		SessionId: sessionId,
		UserId:    userId,
		Name:      name,
	}
	retR, err := cli.RenameSession(context.Background(), req)
	if err != nil {
		log.Error("RenameSession - rpc call:", err.Error())
		ret = false
		return
	}
	ret = retR.Value
	return
}
