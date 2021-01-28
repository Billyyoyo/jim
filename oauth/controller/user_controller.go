package controller

import (
	"github.com/gin-gonic/gin"
	"jim/oauth/dao"
)

//为前端提供的直接的资源服务

//获取用户信息检查accesstoken
func GetUserInfo(c *gin.Context) {
	openId := c.Query("open_id")
	if openId == "" {
		RespError(c, CODE_API_LOSS_OPEN_ID)
		return
	}
	user, err := dao.GetUserByOpenId(openId)
	if err != nil {
		RespError(c, CODE_INTERNAL, err.Error())
		return
	}
	result := map[string]interface{}{
		"open_id":  user.OpenId,
		"nickname": user.NickName,
		"face":     user.Face,
	}
	RespData(c, result)
}
