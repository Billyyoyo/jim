package tests

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"jim/logic/cache"
	"jim/logic/model"
	"testing"
)

func TestSequence(t *testing.T) {
	seq, err := cache.GetUserMsgSequence(1)
	if err != nil {
		log.Error(err.Error())
		return
	}
	fmt.Println(seq)
}

func TestSaveUserConn(t *testing.T) {
	conn1 := &model.UserConn{
		Server:   "localhost:8080",
		Addr:     "127.0.0.1:43331",
		DeviceId: 1,
		Sequence:  3,
	}
	cache.SaveUserConn(1, conn1)
	conn2 := &model.UserConn{
		Server:   "localhost:8082",
		Addr:     "127.0.0.1:43332",
		DeviceId: 1,
		Sequence:  2,
	}
	cache.SaveUserConn(1, conn2)
}

func TestGetUserConn(t *testing.T) {
	conn := &model.UserConn{}
	err := cache.GetUserConn(1, "127.0.0.1:33331", conn)
	if err != nil {
		log.Error(err.Error())
		return
	}
	print(conn)
}

func TestGetAllUserConn(t *testing.T) {
	conns, err := cache.ListUserConn(1)
	if err != nil {
		log.Error(err.Error())
		return
	}
	print(conns)
}

func print(data interface{}) {
	bs, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		log.Error(err.Error())
		return
	}
	fmt.Println(string(bs))
}
