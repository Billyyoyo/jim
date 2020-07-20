package dao

import (
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	log "github.com/sirupsen/logrus"
	"jim/common/utils"
	"jim/logic/core"
	"jim/logic/model"
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

func GetUser(id int64) (user *model.User, err error) {
	user = &model.User{}
	ret, err := db.Where("id=?", id).Get(user)
	if err != nil {
		return
	}
	if !ret {
		err = errors.New("no user")
		return
	}
	return
}

func GetSessionByUser(userId int64) (sessions *[]model.Session, err error) {
	sessions = &[]model.Session{}
	err = db.Table("session").
		Select("session.*").
		Join("INNER", "member", "session.id=member.session_id").
		Where("member.user_id=?", userId).
		Find(sessions)
	return
}

func GetMemberInSession(sessionId int64) (users *[]model.User, err error) {
	users = &[]model.User{}
	err = db.Table("user").
		Select("user.id, user.name").
		Join("INNER", "member", "user.id=member.user_id").
		Where("member.session_id=?", sessionId).
		Find(users)
	return
}

func GetReceptorIdsInSendNo(sessionId, sendNo int64) (ids *[]int64, err error) {
	ids= &[]int64{}
	err= db.Table("message").
		Where("session_id=?", sessionId).
		And("send_no=?",sendNo).
		Cols("receptor_id").
		Find(ids)
	return
}

func IsUserInSession(userId, sessionId int64) (ret bool, err error) {
	ret, err = db.Table("member").
		Where("user_id=?", userId).
		And("session_id=?", sessionId).
		Exist()
	return
}

func AddMember(member *model.Member) (err error) {
	exist, err := db.Table("member").Where("user_id=?", member.UserId).And("session_id=?", member.SessionId).Exist()
	if err != nil {
		return
	}
	if exist {
		err = errors.New("already is member")
		return
	}
	_, err = db.Insert(member)
	if err != nil {
		return
	}
	return
}

func AddMemberV2(ss *xorm.Session, member *model.Member) (err error) {
	exist, err := ss.Table("member").Where("user_id=?", member.UserId).And("session_id=?", member.SessionId).Exist()
	if err != nil {
		return
	}
	if exist {
		err = errors.New("already is member")
		return
	}
	_, err = ss.Insert(member)
	if err != nil {
		return
	}
	return
}

func GetDeviceById(id int64) (device *model.Device, err error) {
	device = &model.Device{}
	ok, err := db.Id(id).Get(device)
	if err != nil {
		return
	}
	if !ok {
		err = errors.New("no device")
		return
	}
	return
}

func GetDevice(userId int64, serial string) (ok bool, device *model.Device, err error) {
	device = &model.Device{}
	ok, err = db.Where("user_id=?", userId).And("serial_no=?", serial).Get(device)
	return
}

func SaveDevice(device *model.Device) (err error) {
	if device.Id > 0 {
		_, err = db.Update(device)
		if err != nil {
			device.Id = 0
			return
		}
	} else {
		_, err = db.Insert(device)
		if err != nil {
			return
		}
	}
	return
}

func SaveDeviceV2(ss *xorm.Session, device *model.Device) (err error) {
	if device.Id > 0 {
		_, err = ss.Update(device)
		if err != nil {
			device.Id = 0
			return
		}
	} else {
		_, err = ss.Insert(device)
		if err != nil {
			return
		}
	}
	return
}

func CreateSession(ss *xorm.Session, session *model.Session) (err error) {
	_, err = ss.Insert(session)
	return
}

func GetSession(sessionId int64) (session *model.Session, err error) {
	session = &model.Session{}
	exist, err := db.Id(sessionId).Get(session)
	if err != nil {
		return
	}
	if !exist {
		err = errors.New("no session")
		return
	}
	return
}

func AddAck(ack *model.Ack) (err error) {
	_, err = db.Insert(ack)
	return
}

func AccumulateAckSendCount(messageId int64) (err error) {
	sql := "update ack set send_count=send_count+1 where msg_id=?"
	_, err = db.Exec(sql, messageId)
	return
}

func AccumulateAckArriveCount(messageId int64) (err error) {
	sql := "update ack set arrive_count=arrive_count+1 where msg_id=?"
	_, err = db.Exec(sql, messageId)
	return
}

func AddMessage(msg *model.Message) (err error) {
	_, err = db.Insert(msg)
	return
}

func GetSessionIdByMessageId(msgId int64) (sessionId int64, err error) {
	sessionId = 0
	exist, err := db.Table(&model.Message{}).ID(msgId).Cols("session_id").Get(&sessionId)
	if err != nil {
		return
	}
	if !exist {
		err = errors.New("no message id")
		return
	}
	return
}

func GetMessagesSeqAfter(receptorId int64, seq int64) (msgs *[]model.Message, err error) {
	msgs = &[]model.Message{}
	err = db.Table(&model.Message{}).
		Where("receptor_id=?", receptorId).
		Where("sequence>?", seq).
		OrderBy("sequence").
		Find(msgs)
	return
}

func GetMessagesSeqIn(receptorId int64, seqs []int64) (msgs *[]model.Message, err error) {
	msgs = &[]model.Message{}
	err = db.Table(&model.Message{}).
		Where("receptor_id=?", receptorId).
		In("sequence", seqs).
		Find(msgs)
	return
}

func GetMessagesSeqRange(receptorId int64, start int64, end int64) (msgs *[]model.Message, err error) {
	msgs = &[]model.Message{}
	err = db.Table(&model.Message{}).
		Where("receptor_id=?", receptorId).
		Where("sequence>?", start).
		And("sequence<?", end).
		Find(msgs)
	return
}

func DeleteMember(sessionId, userId int64) (err error) {
	member := &model.Member{}
	yes, err := db.Table(member).
		Where("session_id=?", sessionId).
		And("user_id=?", userId).
		Get(member)
	if err != nil {
		return
	}
	if !yes {
		err = errors.New("no member")
		return
	}
	_, err = db.Id(member.Id).Delete(member)
	return
}

func RenameSession(sessionId int64, name string) (err error) {
	var session model.Session
	yes, err := db.Id(sessionId).Get(&session)
	if err != nil {
		return
	}
	if !yes {
		return errors.New("no session")
	}
	session.Name = name
	_, err = db.Id(sessionId).Update(session)
	return
}

func WithdrawMessage(userId, sendNo int64) (affect int64, err error) {
	msg := &model.Message{
		Status: 2,
	}
	affect, err = db.
		Where("sender_id=?", userId).
		And("create_time>?", utils.GetCurrentMS()-60_000).
		And("send_no=?", sendNo).
		And("status=?", model.MESSAGE_STATUS_NORMAL).
		Update(msg)
	return
}

func DB() *xorm.Session {
	return db.NewSession()
}
