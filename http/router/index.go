package router

import (
	"github.com/gin-gonic/gin"
	"jim/common/tool"
	"jim/http/controller"
	"strconv"
)

func Route(engine *gin.Engine) {
	apiOneGroup := engine.Group("/jim/api/v1")
	{
		apiOneGroup.GET("/enter", controller.Enter)
	}
	{
		sessionGroup := apiOneGroup.Group("/session")
		sessionGroup.Use(AuthorizationFilter())
		sessionGroup.POST("/create", controller.CreateSession)
		sessionGroup.GET("/list", controller.GetSessions)
		sessionGroup.GET("/get", controller.GetSession)
		sessionGroup.GET("/members", controller.GetMembers)
		sessionGroup.GET("/messages", controller.GetMessages)
		sessionGroup.POST("/rename", controller.RenameSession)
		sessionGroup.POST("/join", controller.JoinSession)
		sessionGroup.POST("/quit", controller.QuitSession)
	}
}

func AuthorizationFilter() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 校验token
		uid := c.GetHeader("jim-uid")
		token := c.GetHeader("jim-token")
		deviceId := c.GetHeader("jim-device")
		if !tool.Validate(uid, deviceId, token) {
			controller.ReturnErr(c, controller.CODE_AUTH)
			c.Abort()
			return
		}
		// 将用户id保存到context
		uidint, _ := strconv.ParseInt(uid, 10, 64)
		controller.SetUserId(c, uidint)
		c.Next()
	}
}
