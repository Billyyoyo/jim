package router

import (
	"github.com/gin-gonic/gin"
	"jim/http/controller"
	"jim/http/dao"
)

func Route(engine *gin.Engine) {
	engine.LoadHTMLGlob("http/templates/*")
	engine.GET("/index", controller.Index)
	engine.GET("/auth/login", controller.Login)
	engine.GET("/auth/callback", controller.AuthCallback)
	apiOneGroup := engine.Group("/jim/api/v1")
	{
		userGroup := apiOneGroup.Group("/user")
		userGroup.Use(AuthorizationFilter())
		userGroup.GET("/self", controller.UserSelf)
	}
	{
		tcpGroup := apiOneGroup.Group("/conn")
		tcpGroup.Use(AuthorizationFilter())
		tcpGroup.GET("/endpoint", controller.GetConnEndPoint)
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
		token := c.GetHeader("jim_token")
		tokenInfo, err := dao.HasToken(token)
		if err != nil {
			controller.ReturnErr(c, controller.CODE_AUTH)
			c.Abort()
			return
		}
		// 将用户id保存到context
		c.Set("jim_user_id", tokenInfo.UserId)
		c.Set("jim_serial_no", tokenInfo.SerialNo)
		c.Set("jim_token", tokenInfo.Token)
		c.Next()
	}
}
