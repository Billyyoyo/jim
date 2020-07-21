package controller

import (
	"github.com/gin-gonic/gin"
	"jim/http/rpc"
)

// 用户请求连接到接入服务器
// 献上用户auth令牌进行验证
// 通过服务发现找到一个连接数最少的接入服务器
// 返回一个接入令牌code
func Enter(c *gin.Context) {

}

//
func GetSessions(c *gin.Context) {

}

//
func GetSession(c *gin.Context) {

}

//
func GetMembers(c *gin.Context) {

}

//
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
	session, err := rpc.CreateSession(GetUserId(c), form.Type, form.Name, form.UserIds)
	if err != nil {
		ReturnError(c, CODE_RPC, err.Error())
		return
	}
	ReturnData(c, session)
}

//
func JoinSession(c *gin.Context) {

}

//
func QuitSession(c *gin.Context) {

}

//
func RenameSession(c *gin.Context) {

}
