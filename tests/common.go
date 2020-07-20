package tests

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
)

func print(data interface{}) {
	bs, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		log.Error(err.Error())
		return
	}
	fmt.Println(string(bs))
}

func println(args ...interface{}) {
	fmt.Println(args...)
}