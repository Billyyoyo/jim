package core

import (
	"encoding/json"
	"errors"
	"fmt"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"time"
)

type appConfig struct {
	Socket struct {
		Mode    string `yaml:"mode"`
		Host    string `yaml:"host"`
		Port    int    `yaml:"port"`
		LogFile string `yaml:"logfile"`
		Tick    int64  `yaml:"tick"`
		Timeout int64  `yaml:"timeout"`
	}
	Rpc struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	}
}

var (
	AppConfig *appConfig
)

func init() {
	confConfig()
	confLog()
}

func confLog() {
	log.Info("------config log-----")
	if AppConfig.Socket.Mode == "release" {
		writer, _ := rotatelogs.New(
			AppConfig.Socket.LogFile+"%Y%m%d.log",
			rotatelogs.WithLinkName(AppConfig.Socket.LogFile),
			rotatelogs.WithMaxAge(time.Hour*24*7),
			rotatelogs.WithRotationTime(time.Hour),
		)
		log.SetOutput(writer)
		log.SetLevel(log.InfoLevel)
	}
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
}

func confConfig() {
	log.Info("------config file-----")
	bs, err := ioutil.ReadFile("/home/billyyoyo/workspace/jim/jim_server/conf/tcp.yaml")
	if err != nil {
		panic(errors.New("config file load error: " + err.Error()))
		return
	}
	AppConfig = &appConfig{}
	err = yaml.Unmarshal(bs, AppConfig)
	if err != nil {
		panic(errors.New("config file parse error: " + err.Error()))
		return
	}
	bs, _ = json.Marshal(AppConfig)
	fmt.Println(string(bs))
}
