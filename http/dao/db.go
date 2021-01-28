package dao

import (
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	log "github.com/sirupsen/logrus"
	"jim/http/core"
	"jim/http/model"
)

var (
	// mysql
	db *xorm.Engine
)

func init() {
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

func GetUser(id int64) (user model.User, err error) {
	user = model.User{}
	ret, err := db.Where("id=?", id).Get(&user)
	if err != nil {
		return
	}
	if !ret {
		err = errors.New("no user")
		return
	}
	return
}

func GetUserByOpenId(openId string) (user model.User, err error) {
	user = model.User{}
	ret, err := db.Where("open_id=?", openId).Get(&user)
	if err != nil {
		return
	}
	if !ret {
		err = errors.New("no user")
		return
	}
	return
}

func SaveUser(user model.User) (id int64, err error) {
	if user.Id > 0 {
		id = user.Id
		_, err = db.Id(user.Id).Update(user)
		if err != nil {
			return
		}
	} else {
		id, err = db.Insert(user)
		if err != nil {
			return
		}
	}
	return
}