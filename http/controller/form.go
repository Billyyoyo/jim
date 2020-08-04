package controller

type CreateSessionForm struct {
	Name    string  `json:"name"`
	Type    int8    `json:"type"`
	UserIds []int64 `json:"userIds"`
}

type SessionActionForm struct {
	UserId    int64  `json:"userId"`
	SessionId int64  `json:"sessionId"`
	Name      string `json:"name"`
}

type GetMessagesForm struct {
	SessionId int64  `form:"sessionId"`
	Cond      string `form:"cond"`
}
