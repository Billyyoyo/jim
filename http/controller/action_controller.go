package controller

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"jim/common/tool"
)

// 用户请求连接到接入服务器
// 返回一个接入token
func Enter(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		ReturnErr(c, CODE_PARAMS)
		return
	}
	// 献上用户auth令牌进行验证
	ok, uid, token := tool.Authorization(code)
	if !ok {
		ReturnError(c, CODE_RPC, "授权失败")
		return
	}
	// 通过服务发现找到一个连接数最少的接入服务器 todo 返回一个固定接入服务器地址
	address := "localhost:4002"

	result := map[string]interface{}{
		"tcpServer": address,
		"uid":       uid,
		"token":     token,
	}
	ReturnData(c, result)
}

// 获取用户所有会话
func GetSessions(c *gin.Context) {
	sessions, err := tool.GetSessions(Uid(c))
	if err != nil {
		log.Error("GetSessions fail:", Uid(c))
		ReturnErr(c, CODE_RPC)
		return
	}
	ReturnData(c, sessions)
}

// 获取会话及成员
func GetSession(c *gin.Context) {
	sessionId, err := QueryInt(c, "sessionId")
	if err != nil {
		ReturnErr(c, CODE_PARAMS)
		return
	}
	session, err := tool.GetSession(sessionId)
	if err != nil {
		log.Error("GetSession fail:", Uid(c))
		ReturnErr(c, CODE_RPC)
		return
	}
	ReturnData(c, session)
}

// 获取会话中所有成员
func GetMembers(c *gin.Context) {
	sessionId, err := QueryInt(c, "sessionId")
	if err != nil {
		ReturnErr(c, CODE_PARAMS)
		return
	}
	members, err := tool.GetMembers(sessionId)
	if err != nil {
		log.Error("GetMembers fail:", sessionId)
		ReturnErr(c, CODE_RPC)
		return
	}
	ReturnData(c, members)
}

// 创建会话
func CreateSession(c *gin.Context) {
	form := &CreateSessionForm{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		ReturnError(c, CODE_PARAMS, err.Error())
		return
	}
	if form.Name == "" {
		ReturnError(c, CODE_PARAMS, "名称不能为空")
		return
	}
	if form.Type != 1 && form.Type != 2 {
		ReturnError(c, CODE_PARAMS, "类型错误")
		return
	}
	if len(form.UserIds) < 2 {
		ReturnError(c, CODE_PARAMS, "人数错误")
		return
	}
	session, err := tool.CreateSession(Uid(c), form.Type, form.Name, form.UserIds)
	if err != nil {
		ReturnError(c, CODE_RPC, err.Error())
		log.Error("Create session - rpc faild")
		return
	}
	ReturnData(c, session)
}

// 加入会话
func JoinSession(c *gin.Context) {
	form := &SessionActionForm{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		ReturnError(c, CODE_PARAMS, err.Error())
		return
	}
	if form.UserId == 0 {
		ReturnError(c, CODE_PARAMS, "用户不能为空")
		return
	}
	if form.SessionId == 0 {
		ReturnError(c, CODE_PARAMS, "群组不能为空")
		return
	}
	ret, err := tool.JoinSession(form.UserId, form.SessionId)
	if err != nil || !ret {
		ReturnError(c, CODE_RPC, "不能加入群组")
		log.Error("Join session - rpc faild")
		return
	}
	ReturnSuccess(c)
}

// 退出会话
func QuitSession(c *gin.Context) {
	form := &SessionActionForm{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		ReturnError(c, CODE_PARAMS, err.Error())
		return
	}
	if form.UserId == 0 {
		ReturnError(c, CODE_PARAMS, "用户不能为空")
		return
	}
	if form.SessionId == 0 {
		ReturnError(c, CODE_PARAMS, "群组不能为空")
		return
	}
	ret, err := tool.QuitSession(form.UserId, form.SessionId)
	if err != nil || !ret {
		ReturnError(c, CODE_RPC, "操作失败")
		log.Error("Quit session - rpc faild")
		return
	}
	ReturnSuccess(c)
}

// 重命名会话
func RenameSession(c *gin.Context) {
	form := &SessionActionForm{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		ReturnError(c, CODE_PARAMS, err.Error())
		return
	}
	if form.UserId == 0 {
		ReturnError(c, CODE_PARAMS, "用户不能为空")
		return
	}
	if form.SessionId == 0 {
		ReturnError(c, CODE_PARAMS, "群组不能为空")
		return
	}
	if form.Name == "" {
		ReturnError(c, CODE_PARAMS, "群名称不能为空")
		return
	}
	ret, err := tool.RenameSession(form.UserId, form.SessionId, form.Name)
	if err != nil || !ret {
		log.Error("Rename session - rpc faild")
		ReturnError(c, CODE_RPC, "操作失败")
		return
	}
	ReturnSuccess(c)
}
