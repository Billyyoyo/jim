package tests

import (
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"jim/common/rpc"
)

var (
	cli  rpc.LogicServiceClient
	cli2 rpc.SocketServiceClient
)

func init() {
	conn, err := grpc.Dial("localhost:4000", grpc.WithInsecure())
	if err != nil {
		panic("grpc start up error: " + err.Error())
		return
	}
	cli = rpc.NewLogicServiceClient(conn)

	conn2, err := grpc.Dial("localhost:4003", grpc.WithInsecure())
	if err != nil {
		panic("grpc start up error: " + err.Error())
		return
	}
	cli2 = rpc.NewSocketServiceClient(conn2)
}

func printj(data interface{}) {
	bs, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		log.Error(err.Error())
		return
	}
	fmt.Println(string(bs))
}

func printl(args ...interface{}) {
	fmt.Println(args...)
}

func printp(data proto.Message) {
	bs, err := proto.Marshal(data)
	if err != nil {
		log.Error(err.Error())
		return
	}
	fmt.Println("{", string(bs), "}")
}
