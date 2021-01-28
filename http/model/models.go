package model

import (
	"encoding/json"
	"fmt"
)

type User struct {
	Id         int64  `xorm:"'id' pk autoincr"`
	Name       string `xorm:"name"`
	OpenId     string `xorm:"open_id"`
	Face       string `xorm:"face"`
	AuthToken  string `xorm:"auth_token"`
	CreateTime int64  `xorm:"create_time"`
}

func (u User) TableName() string {
	return "im_user"
}

type TokenInfo struct {
	Token string
	UserId int64
	SerialNo string
}

func (t *TokenInfo) MarshalBinary() ([]byte, error) {
	return json.Marshal(t)
}

func (t *TokenInfo) UnmarshalBinary(data []byte) error {
	fmt.Println(string(data))
	return json.Unmarshal(data, t)
}

