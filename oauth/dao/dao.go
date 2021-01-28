package dao

import (
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	log "github.com/sirupsen/logrus"
	"jim/oauth/core"
	"jim/oauth/model"
	"time"
)

var (
	cache *redis.Client
	db    *xorm.Engine
)

func init() {
	cache = redis.NewClient(&redis.Options{
		Addr:     core.AppConfig.Redis.Addr,
		Password: "",
		DB:       0,
	})
	var err error
	address := fmt.Sprintf(
		"%s:%s@(%s:%s)/%s?charset=utf8",
		core.AppConfig.Database.User,
		core.AppConfig.Database.Password,
		core.AppConfig.Database.Addr,
		core.AppConfig.Database.Port,
		core.AppConfig.Database.DB)
	log.Info("database: ", "mysql\t", "connection: ", address)
	db, err = xorm.NewEngine("mysql", address)
	if err != nil {
		panic("db init failed:" + err.Error())
		return
	}
	if core.AppConfig.Server.Mode == "debug" {
		db.ShowSQL(true)
	} else {
		db.ShowSQL(false)
	}
}

func GetUserById(id int64) (user model.User, err error) {
	user = model.User{}
	ok, err := db.ID(id).Get(&user)
	if err != nil {
		return
	}
	if !ok {
		err = errors.New(fmt.Sprintf("UserId %d not exist", id))
	}
	return
}

func GetUserByLoginName(loginName string) (user model.User, err error) {
	user = model.User{}
	ok, err := db.Table(user).Where("login_name=?", loginName).Get(&user)
	if err != nil {
		return
	}
	if !ok {
		err = errors.New(fmt.Sprintf("Login name %s not exist", loginName))
	}
	return
}

func GetUserByOpenId(openId string) (user model.User, err error) {
	user = model.User{}
	ok, err := db.Table(user).Where("open_id=?", openId).Get(&user)
	if err != nil {
		return
	}
	if !ok {
		err = errors.New(fmt.Sprintf("OpenId %s not exist", openId))
	}
	return
}

func VerrifyAndGetUser(loginName, password string) (user model.User, err error) {
	user = model.User{}
	ok, err := db.Where("login_name=? and password=?", loginName, password).Get(&user)
	if err != nil {
		return
	}
	if !ok {
		err = errors.New(fmt.Sprintf("LoginName or password is wrong"))
	}
	return
}

func SaveNewUser(user model.User) (userId int64, err error) {
	isExist, err := db.Exist(&model.User{LoginName: user.LoginName})
	if err != nil {
		return
	}
	if isExist {
		err = errors.New("LoginName already exist")
		return
	}
	userId, err = db.Insert(user)
	return
}

func GetUserList() (users []model.User, err error) {
	users = []model.User{}
	err = db.Table(model.User{}).Find(&users)
	return
}

func GetApplicationById(id int64) (app model.Application, err error) {
	app = model.Application{}
	ret, err := db.Id(id).Get(&app)
	if !ret {
		err = errors.New("not exists")
	}
	return
}

func GenerateMainToken(userId int64) string {
	text := fmt.Sprintf("jim_token_%d_%d", userId, time.Now().UnixNano())
	return fmt.Sprintf("%x", md5.Sum([]byte(text)))
}

func GenerateCertificateCode(userId, appId int64) string {
	text := fmt.Sprintf("jim_cert_code_%d_%d_%d", appId, userId, time.Now().UnixNano())
	return fmt.Sprintf("%x", md5.Sum([]byte(text)))
}
