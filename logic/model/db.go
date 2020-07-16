package model

type Message struct {
	Id         int64  `xorm:"'id' pk autoincr"`
	SenderId   int64  `xorm:"sender_id"`
	SessionId  int64  `xorm:"session_id"`
	Type       int8   `xorm:"type"`
	Status     int8   `xorm:"status"`
	DeviceId   int64  `xorm:"device_id"`
	Sequence   int64  `xorm:"sequence"`
	ReceptorId int64  `xorm:"receptor_id"`
	Body       []byte `xorm:"body"`
	CreateTime int64  `xorm:"create_time"`
}

func (m *Message) TableName() string {
	return "message"
}

type User struct {
	Id         int64  `xorm:"'id' pk autoincr"`
	Name       string `xorm:"name"`
	Pwd        string `xorm:"password"`
	CreateTime int64  `xorm:"create_time"`
}

func (u *User) TableName() string {
	return "user"
}

type Session struct {
	Id         int64  `xorm:"'id' pk autoincr"`
	Name       string `xorm:"name"`
	Type       int8   `xorm:"type"`
	Owner      int64  `xorm:"owner"`
	CreateTime int64  `xorm:"create_time"`
}

func (s *Session) TableName() string {
	return "session"
}

type Member struct {
	Id         int64 `xorm:"'id' pk autoincr"`
	SessionId  int64 `xorm:"session_id"`
	UserId     int64 `xorm:"user_id"`
	CreateTime int64 `xorm:"create_time"`
}

func (m *Member) TableName() string {
	return "member"
}

type Device struct {
	Id             int64  `xorm:"'id' pk autoincr"`
	UserId         int64  `xorm:"user_id"`
	SerialNo       string `xorm:"serial_no"`
	LastAddress    string `xorm:"last_address"`
	LastConnTime   int64  `xorm:"last_conn_time"`
	LastDisconTime int64  `xorm:"last_discon_time"`
	LastSequence   int64  `xorm:"last_sequence"`
	CreateTime     int64  `xorm:"create_time"`
}

func (d *Device) TableName() string {
	return "device"
}

type Ack struct {
	Id          int64 `xorm:"'id' pk autoincr"`
	MsgId       int64 `xorm:"msg_id"`
	SendCount   int   `xorm:"send_count"`
	ArriveCount int   `xorm:"arrive_count"`
}

func (a *Ack) TableName() string {
	return "ack"
}
