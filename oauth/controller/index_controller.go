package controller

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"jim/common/utils"
	"jim/oauth/core"
	"jim/oauth/dao"
	"jim/oauth/model"
	"net/http"
	"strconv"
	"strings"
)

//判断cookie里面有没有token及有效性
//true 跳转returnurl
//false 跳转login
func Authorize(c *gin.Context) {
	appId, _ := QueryInt64(c, "app_id")
	extra := c.Query("extra")
	if appId == 0 {
		c.AbortWithError(http.StatusBadRequest, errors.New("miss param app_id"))
		return
	}
	app, err := dao.GetApplicationById(appId)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	token, _ := c.Cookie("j_token")
	if token == "" {
		redirectTo(c, "login")
		return
	}

	userId, err := dao.GetOAuthToken(token)
	if err != nil || userId == 0 {
		c.SetCookie("j_token", "", -1, "/", "localhost", false, false)
		redirectTo(c, "login")
		return
	}

	certCode := dao.GenerateCertificateCode(userId, appId)
	if err = dao.SaveCertCode(userId, appId, certCode); err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	returnUrl := c.Query("return_url")
	idx := strings.Index(returnUrl, app.Host)
	if returnUrl == "" || idx != 0 {
		returnUrl = app.Host
	}
	returnUrl = fmt.Sprintf("%s?code=%s&extra=%s", returnUrl, certCode, extra)
	c.Redirect(http.StatusFound, returnUrl)
}

func Welcome(c *gin.Context) {
	token, _ := c.Cookie("j_token")
	if token == "" {
		redirectTo(c, "login")
		return
	}
	userId, err := dao.GetOAuthToken(token)
	if err != nil || userId == 0 {
		c.SetCookie("j_token", "", -1, "/", "localhost", false, false)
		redirectTo(c, "login")
		return
	}
	view := gin.H{
		"UserId": userId,
	}
	c.HTML(http.StatusOK, "welcome.tmpl", view)
}

func Error(c *gin.Context) {
	view := gin.H{
		"message": c.Query("message"),
	}
	c.HTML(http.StatusOK, "welcome.tmpl", view)
}

//用户注册
//注册成功后，生成授权码
func ToRegister(c *gin.Context) {
	token, _ := c.Cookie("j_token")
	if token != "" {
		redirectTo(c, "authorize")
		return
	}
	views := gin.H{
		"return_url": c.Query("return_url"),
		"app_id":     c.Query("app_id"),
		"extra":     c.Query("extra"),
	}
	c.HTML(http.StatusOK, "register.tmpl", views)
}

//用户登录（页面包含注册按钮）
//登录成功后，生成授权码
func ToLogin(c *gin.Context) {
	token, _ := c.Cookie("j_token")
	if token != "" {
		redirectTo(c, "authorize")
		return
	}
	views := gin.H{
		"return_url": c.Query("return_url"),
		"app_id":     c.Query("app_id"),
		"extra":     c.Query("extra"),
	}
	c.HTML(http.StatusOK, "login.tmpl", views)
}

//用户注册
//注册成功后，生成授权码
func Register(c *gin.Context) {
	loginName := c.PostForm("login_name")
	if loginName == "" {
		c.AbortWithError(http.StatusBadRequest, errors.New("login name is empty"))
		return
	}
	nickname := c.PostForm("nickname")
	if nickname == "" {
		c.AbortWithError(http.StatusBadRequest, errors.New("nickname is empty"))
		return
	}
	password := c.PostForm("password")
	if password == "" {
		c.AbortWithError(http.StatusBadRequest, errors.New("password is empty"))
		return
	}
	password = utils.Md5sum(loginName + password)
	user := model.User{
		OpenId:    strings.ReplaceAll(uuid.New().String(), "-", ""),
		LoginName: loginName,
		NickName:  nickname,
		Password:  password,
	}
	userId, err := dao.SaveNewUser(user)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	user.Id = userId
	success(c, userId)
}

//用户登录（页面包含注册按钮）
//登录成功后，生成授权码
func Login(c *gin.Context) {
	loginName := c.PostForm("login_name")
	if loginName == "" {
		c.AbortWithError(http.StatusBadRequest, errors.New("login name is empty"))
		return
	}
	password := c.PostForm("password")
	if password == "" {
		c.AbortWithError(http.StatusBadRequest, errors.New("password is empty"))
		return
	}
	password = utils.Md5sum(loginName + password)
	user, err := dao.VerrifyAndGetUser(loginName, password)
	if err != nil {
		views := gin.H{
			"return_url": c.Query("return_url"),
			"app_id":     c.Query("app_id"),
			"extra":     c.Query("extra"),
			"err_msg":    err.Error(),
		}
		c.HTML(http.StatusOK, "login.tmpl", views)
		return
	}
	success(c, user.Id)
}

func success(c *gin.Context, userId int64) {
	appId_str := c.PostForm("app_id")
	appId, _ := strconv.ParseInt(appId_str, 10, 64)

	returnUrl := c.PostForm("return_url")
	token := dao.GenerateMainToken(userId)
	err := dao.SaveOAuthToken(userId, token)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.SetCookie("j_token", token, 24*60*60, "/", core.AppConfig.Server.Host, false, false)
	if appId == 0 || returnUrl == "" {
		c.Redirect(http.StatusFound, "welcome")
		return
	}
	app, err := dao.GetApplicationById(appId)
	if err != nil {
		c.Redirect(http.StatusFound, "error?message=找不到应用")
		return
	}
	if strings.Index(returnUrl, app.Host) != 0 {
		returnUrl = app.Host
	}
	certCode := dao.GenerateCertificateCode(userId, appId)
	if err = dao.SaveCertCode(userId, appId, certCode); err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}
	extra := c.PostForm("extra")
	returnUrl = fmt.Sprintf("%s?code=%s&extra=%s", returnUrl, certCode, extra)
	c.Redirect(http.StatusFound, returnUrl)
}

func redirectTo(c *gin.Context, path string) {
	c.Redirect(http.StatusFound, fmt.Sprintf("%s://%s:%d/%s?app_id=%s&return_url=%s&extra=%s",
		core.AppConfig.Server.Schema,
		core.AppConfig.Server.Host,
		core.AppConfig.Server.Port,
		path,
		c.Query("app_id"),
		c.Query("return_url"),
		c.Query("extra")))
}

//用户授权
//起点接口，检查请求中是否有appid和跳转地址，如没有返回错误界面，检查是否有授权码，如没有跳转登录页面(参数带有appid和跳转地址)
//生成accesstoken和refreshtoken保存服务端
//生成登录码返回给调用的前端
//改接口需要验签
func Certificate(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		RespError(c, CODE_AUTH_CODE_ERROR, "miss auth code")
		return
	}
	userId, appId := dao.GetCertCode(code)
	if userId == 0 {
		RespError(c, CODE_AUTH_CODE_ERROR)
		return
	}
	user, err := dao.GetUserById(userId)
	if err != nil {
		RespError(c, CODE_AUTH_CODE_ERROR, "no user found")
		return
	}
	authToken, err := dao.GenerateAuthToken(userId, appId)
	result := map[string]interface{}{
		"open_id":       user.OpenId,
		"auth_token":  authToken,
	}
	RespData(c, result)
}

