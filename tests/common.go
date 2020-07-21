package tests

import (
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
)

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
