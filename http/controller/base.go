package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"jim/common"
)

const (
	CODE_SUCCESS = iota
	CODE_AUTH
	CODE_PERMISSION
	CODE_INTERNAL
	CODE_RPC
	CODE_PARAMS
	CODE_NO_USER_INFO
	CODE_RESTFUL
)

var (
	CODES_NAME = []string{
		"ok",
		"未登录",
		"没有权限",
		"内部错误",
		"远程调用错误",
		"参数错误",
		"没有用户信息",
		"服务调用失败",
	}
)

func ReturnSuccess(c *gin.Context) {
	c.JSON(http.StatusOK, common.Result{
		Code: CODE_SUCCESS,
		Msg:  CODES_NAME[CODE_SUCCESS],
	})
}

func ReturnData(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, common.Result{
		Code: CODE_SUCCESS,
		Msg:  CODES_NAME[CODE_SUCCESS],
		Data: data,
	})
}

func ReturnErr(c *gin.Context, code int) {
	var httpStatus int
	if code == CODE_AUTH {
		httpStatus = http.StatusUnauthorized
	} else if code == CODE_PERMISSION {
		httpStatus = http.StatusForbidden
	} else {
		httpStatus = http.StatusOK
	}
	c.JSON(httpStatus, common.Result{
		Code: code,
		Msg:  CODES_NAME[code],
		Desc: CODES_NAME[code],
	})
}

func ReturnError(c *gin.Context, code int, desc string) {
	var httpStatus int
	if code == CODE_AUTH {
		httpStatus = http.StatusUnauthorized
	} else if code == CODE_PERMISSION {
		httpStatus = http.StatusForbidden
	} else {
		httpStatus = http.StatusOK
	}
	c.JSON(httpStatus, common.Result{
		Code: code,
		Msg:  CODES_NAME[code],
		Desc: desc,
	})
}

func Uid(c *gin.Context) int64 {
	return c.GetInt64("jim_user_id")
}

func Utoken(c *gin.Context) string {
	return c.GetString("jim_token")
}

func Userial(c *gin.Context) string {
	return c.GetString("jim_serial_no")
}

func QueryInt(c *gin.Context, key string) (value int64, err error) {
	v := c.Query(key)
	return strconv.ParseInt(v, 10, 64)
}
