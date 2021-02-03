package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"time"
)

type appConfig struct {
	Server struct {
		Mode    string `yaml:"mode"`
		Host    string `yaml:"host"`
		Port    int    `yaml: "port"`
		LogFile string `yaml: "logfile"`
	}

	Redis struct {
		Addr   string `yaml:"addr"`
		Prefix string `yaml: "prefix"`
	}

	Database struct {
		Addr     string `yaml:"addr"`
		Port     string `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		DB       string `yaml:"db"`
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
	if AppConfig.Server.Mode == "release" {
		writer, _ := rotatelogs.New(
			AppConfig.Server.LogFile+"%Y%m%d.log",
			rotatelogs.WithLinkName(AppConfig.Server.LogFile),
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
	bs, err := ioutil.ReadFile("/home/billyyoyo/workspace/jim/jim_server/conf/http.yaml")
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

func Logger() gin.HandlerFunc {
	logger := log.StandardLogger()
	return func(c *gin.Context) {
		startTime := time.Now()
		c.Next()
		endTime := time.Now()
		latencyTime := endTime.Sub(startTime)
		reqMethod := c.Request.Method
		reqUri := c.Request.RequestURI
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		logger.Infof("| %3d | %13v | %15s | %s | %s |",
			statusCode,
			latencyTime,
			clientIP,
			reqMethod,
			reqUri,
		)
	}
}

func CROS() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "jim_token")
		c.Header("Access-Control-Allow-Methods", "POST, PUT, DELETE, GET, OPTIONS")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}
