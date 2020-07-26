package model

import "encoding/json"

// todo 根据redis特点    这里还有优化空间
type UserState struct {
	Server   string
	Addr     string
	DeviceId int64
	Token    string
}

func (u *UserState) MarshalBinary() ([]byte, error) {
	return json.Marshal(u)
}

func (u *UserState) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, u)
}
