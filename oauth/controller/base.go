package controller

import (
	"github.com/gin-gonic/gin"
	"jim/common"
	"net/http"
	"strconv"
)

const (
	CODE_SUCCESS = iota
	CODE_INTERNAL
	CODE_AUTH_CODE_ERROR
	CODE_AUTH_AUTH_TOKEN_ERROR
	CODE_AUTH_LOSS_APP_ID
	CODE_API_LOSS_OPEN_ID
	CODE_API_LOSS_ACCESS_TOKEN
)

var CODE_MSG = []string{
	"success",
	"internal error",
	"Authorization code is wrong",
	"Auth token is wrong",
	"loss appid in parameters",
	"loss openid in parameters",
	"api loss access token",
}

func QueryInt64(c *gin.Context, key string) (value int64, err error) {
	v := c.Query(key)
	return strconv.ParseInt(v, 10, 64)
}

func RespSuccess(c *gin.Context) {
	c.JSON(http.StatusOK, common.Result{
		Code: CODE_SUCCESS,
		Msg:  CODE_MSG[CODE_SUCCESS],
	})
}

func RespData(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, common.Result{
		Code: CODE_SUCCESS,
		Msg:  CODE_MSG[CODE_SUCCESS],
		Data: data,
	})
}

func RespError(c *gin.Context, code int, desc ...string) {
	result := common.Result{
		Code: code,
		Msg:  CODE_MSG[code],
	}
	if len(desc) > 0 {
		result.Desc = desc[0]
	}
	c.JSON(http.StatusOK, result)
}
