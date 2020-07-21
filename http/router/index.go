package router

import (
	"github.com/gin-gonic/gin"
	"jim/http/controller"
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
		sessionGroup.GET("/list", controller.CreateSession)
		sessionGroup.GET("/get", controller.GetSession)
		sessionGroup.GET("/members", controller.GetMembers)
		sessionGroup.POST("/rename", controller.RenameSession)
		sessionGroup.POST("/join", controller.JoinSession)
		sessionGroup.POST("/quit", controller.QuitSession)
	}
}

func AuthorizationFilter() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 校验token
		token := c.GetHeader("jim-token")
		if token == "" {
			controller.ReturnErr(c, controller.CODE_AUTH)
			c.Abort()
			return
		}
		c.Next()
	}
}
