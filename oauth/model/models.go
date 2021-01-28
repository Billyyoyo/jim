package model

type User struct {
	Id        int64  `xorm:"'id' pk autoincr"`
	OpenId    string `xorm:"open_id"`
	LoginName string `xorm:"login_name"`
	NickName  string `xorm:"nickname"`
	Face      string `xorm:"face"`
	Password  string `xorm:"password"`
	Flag      int8   `xorm:"'flag'"`
}

func (m User) TableName() string {
	return "auth_user"
}

type Application struct {
	Id        int64  `xorm:"'id' pk autoincr"`
	Name      string `xorm:"name"`
	Key       string `xorm:"key"`
	PublicKey string `xorm:"public_key"`
	Host      string `xorm:"host"`
}

func (m Application) TableName() string {
	return "auth_application"
}

type AuthCode struct {
	OpenId        string `json:"open_id"`
	Expired       int64  `json:"expired_time"`
	ApplicationId int64  `json:"application_id"`
}
