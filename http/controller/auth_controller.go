package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"jim/common/utils"
	"jim/http/core"
	"jim/http/dao"
	"jim/http/model"
	"jim/http/service"
	"net/http"
	"net/url"
	"time"
)

func Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl", gin.H{})
}

func Login(c *gin.Context) {
	serialNo := c.Query("serial_no")
	if serialNo == "" {
		ReturnError(c, CODE_PARAMS, "no device serial no")
		return
	}
	appToken, err := c.Cookie("jim_token")
	if err != nil || appToken == "" {
		goOauth(c, serialNo)
		return
	}
	tokenInfo, err := dao.HasToken(appToken)
	if err != nil {
		goOauth(c, serialNo)
		return
	}
	err = dao.RefreshToken(tokenInfo.Token)
	if err != nil {
		ReturnError(c, CODE_PERSISTENCE, err.Error())
		return
	}
	c.SetCookie("jim_token", appToken, 7*24*60*60, "/", core.AppConfig.Server.Host, false, false)
	c.HTML(http.StatusOK, "auth_success.tmpl", gin.H{})
}

func goOauth(c *gin.Context, serialNo string) {
	returnUrl := url.QueryEscape(core.AppConfig.Server.OauthCallback)
	url := fmt.Sprintf("%s?app_id=%s&return_url=%s&extra=%s",
		core.AppConfig.Server.OauthApi,
		core.AppConfig.Server.OauthAppId,
		returnUrl,
		serialNo)
	c.Redirect(http.StatusFound, url)
}

type UserResp struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
	Face string `json:"face"`
}

func AuthCallback(c *gin.Context) {
	code := c.Query("code")
	serialNo := c.Query("extra")
	if code == "" {
		ReturnErr(c, CODE_PARAMS)
		return
	}
	token, er := service.Certificate(code)
	if er != nil {
		ReturnErr(c, CODE_AUTH)
		return
	}
	user, err := dao.GetUserByOpenId(token.OpenId)
	if err != nil {
		user = model.User{
			OpenId:     token.OpenId,
			CreateTime: utils.GetCurrentMS(),
		}
		userInfo, errr := service.GetUserInfo(token.OpenId, token.AuthToken)
		if errr != nil {
			ReturnErr(c, CODE_NO_USER_INFO)
			return
		}
		user.Name = userInfo.Nickname
		user.Face = userInfo.Face
	}
	user.AuthToken = token.AuthToken
	userId, err := dao.SaveUser(user)
	if err != nil {
		ReturnErr(c, CODE_INTERNAL)
		return
	}
	user.Id = userId
	appToken := utils.Md5sum(fmt.Sprintf("app_token_%d_%d", user.Id, time.Now().Nanosecond()))
	err = dao.SaveToken(user.Id, serialNo, appToken)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.SetCookie("jim_token", appToken, 7*24*60*60, "/", core.AppConfig.Server.Host, false, false)
	c.HTML(http.StatusOK, "auth_success.tmpl", gin.H{})
}
