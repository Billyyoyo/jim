package router

import "C"
import (
	"github.com/gin-gonic/gin"
	"jim/oauth/controller"
	"jim/oauth/dao"
	"net/http"
)

func Route(engine *gin.Engine) {
	engine.LoadHTMLGlob("oauth/templates/*")
	engine.StaticFS("/oauth/home", http.Dir("./oauth/static"))
	engine.GET("/authorize", controller.Authorize)
	engine.GET("/welcome", controller.Welcome)
	engine.GET("/error", controller.Error)
	engine.GET("/login", controller.ToLogin)
	engine.GET("/register", controller.ToRegister)
	engine.POST("/login", controller.Login)
	engine.POST("/register", controller.Register)
	apiAuthGroup := engine.Group("/oauth/api/v1")
	//auth
	{
		apiAuthGroup.GET("certificate", controller.Certificate)
	}
	//resource-user
	apiResGroup := engine.Group("/resource/api/v1")
	{
		apiResGroup.Use(AuthFilter())
		apiResGroup.GET("user/info", controller.GetUserInfo)
	}
}

func AuthFilter() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 校验token
		token := c.Query("auth_token")
		if token == "" {
			controller.RespError(c, controller.CODE_API_LOSS_ACCESS_TOKEN)
			c.Abort()
			return
		}
		userId, _, err := dao.ValidateAuthToken(token)
		if err != nil {
			controller.RespError(c, controller.CODE_AUTH_AUTH_TOKEN_ERROR)
			c.Abort()
			return
		}
		// 将用户id保存到context
		c.Set("user_id", userId)
		c.Next()
	}
}
