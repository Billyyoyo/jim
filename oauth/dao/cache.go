package dao

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"jim/common/utils"
	"jim/oauth/core"
	"strconv"
	"strings"
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

func GetOAuthToken(token string) (userId int64, err error) {
	cmd := client.Get(context.Background(), fmt.Sprintf("oauth_token_%s", token))
	userId, err = cmd.Int64()
	return
}

func SaveOAuthToken(userId int64, token string) (err error) {
	expiredTime := time.Duration(core.AppConfig.Server.MainTokenExpired) * 24 * time.Hour
	cmd := client.Set(context.Background(), fmt.Sprintf("oauth_token_%s", token), userId, expiredTime)
	_, err = cmd.Result()
	return
}

func GetCertCode(code string) (userId, appId int64) {
	key := fmt.Sprintf("oauth_cert_code_%s", code)
	cmd := client.Get(context.Background(), key)
	str, err := cmd.Result()
	if err != nil {
		return
	}
	strs := strings.Split(str, "|")
	userId, _ = strconv.ParseInt(strs[0], 10, 64)
	appId, _ = strconv.ParseInt(strs[1], 10, 64)
	client.Del(context.Background(), key)
	return
}

func SaveCertCode(userId, appId int64, code string) (err error) {
	expiredTime := 5 * time.Minute
	value := fmt.Sprintf("%d|%d", userId, appId)
	cmd := client.Set(context.Background(), fmt.Sprintf("oauth_cert_code_%s", code), value, expiredTime)
	_, err = cmd.Result()
	return
}

func GenerateAuthToken(userId, appId int64) (authToken string, err error) {
	accessExpired := time.Duration(core.AppConfig.Server.AuthTokenExpired) * 24 * time.Hour
	flag := fmt.Sprintf("%d|%d", appId, userId)
	authToken = genToken("jt_auth_", flag)
	client.Set(context.Background(), "j_at_"+authToken, flag, accessExpired)
	return
}

func ValidateAuthToken(accessToken string) (userId, appId int64, err error) {
	cmd := client.Get(context.Background(), "j_at_"+accessToken)
	str, err := cmd.Result()
	if err != nil {
		return
	}
	strs := strings.Split(str, "|")
	userId, _ = strconv.ParseInt(strs[0], 10, 64)
	appId, _ = strconv.ParseInt(strs[1], 10, 64)
	client.Set(context.Background(), "j_at_"+accessToken, str,
		time.Duration(core.AppConfig.Server.AuthTokenExpired) * 24 * time.Hour)
	return
}

func genToken(prex, flag string) string {
	block1 := utils.Md5sum(fmt.Sprintf("%s%s_%d", prex, flag, time.Now().UnixNano()))
	block2 := utils.EncryptDES([]byte(flag), []byte(utils.DES_KEY))
	return fmt.Sprintf("%s.%s", block1, block2)
}
