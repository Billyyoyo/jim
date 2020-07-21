package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	CODE_SUCCESS = iota
	CODE_AUTH
	CODE_PERMISSION
	CODE_INTERNAL
	CODE_RPC
	CODE_PARAMS
)

var (
	CODES_NAME = []string{
		"ok",
		"未登录",
		"没有权限",
		"内部错误",
		"远程调用错误",
		"参数错误",
	}
)

type Result struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Desc string      `json:"desc"`
	Data interface{} `json:"data"`
}

func ReturnSuccess(c *gin.Context) {
	c.JSON(http.StatusOK, Result{
		Code: CODE_SUCCESS,
		Msg:  CODES_NAME[CODE_SUCCESS],
	})
}

func ReturnData(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Result{
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
	c.JSON(httpStatus, Result{
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
	c.JSON(httpStatus, Result{
		Code: code,
		Msg:  CODES_NAME[code],
		Desc: desc,
	})
}

func SetUserId(c *gin.Context, userId int64) {
	c.Set("jim_user_id", userId)
}

func GetUserId(c *gin.Context) int64 {
	return c.GetInt64("jim_user_id")
}
