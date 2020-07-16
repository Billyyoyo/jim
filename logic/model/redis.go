package model

import "encoding/json"

type UserConn struct {
	Server   string
	UserId   int64
	Addr     string
	DeviceId int64
}

func (u *UserConn) MarshalBinary() ([]byte, error) {
	return json.Marshal(u)
}

func (u *UserConn) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, u)
}
