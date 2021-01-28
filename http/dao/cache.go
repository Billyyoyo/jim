package dao

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"jim/http/core"
	"jim/http/model"
	"time"
)

var (
	client *redis.Client
)

func init() {
	client = redis.NewClient(&redis.Options{
		Addr:     core.AppConfig.Redis.Addr,
		Password: "",
		DB:       0,
	})
}

func SaveToken(userId int64, serialNo, token string) (err error) {
	key := fmt.Sprintf("%s_user_token_%s", core.AppConfig.Redis.Prefix, token)
	tokenInfo := model.TokenInfo{
		UserId:   userId,
		SerialNo: serialNo,
		Token:    token,
	}
	cmd := client.SetNX(context.Background(), key, &tokenInfo, time.Hour*24*7)
	if cmd.Err() != nil {
		err = cmd.Err()
		return
	}
	return
}

func HasToken(token string) (tokenInfo model.TokenInfo, err error) {
	key := fmt.Sprintf("%s_user_token_%s", core.AppConfig.Redis.Prefix, token)
	cmd := client.Get(context.Background(), key)
	if cmd.Err() != nil {
		err = cmd.Err()
		return
	}
	tokenInfo = model.TokenInfo{}
	err = cmd.Scan(&tokenInfo)
	return
}
