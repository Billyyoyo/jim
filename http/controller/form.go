package controller



type CreateSessionForm struct {
	Name    string  `json:"name"`
	Type    int8    `json:"type"`
	UserIds []int64 `json:"userIds"`
}
