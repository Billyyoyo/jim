package service

import (
	"jim/common"
	tool2 "jim/common/tool"
)

type UserInfoResp struct {
	Nickname string
	Face     string
	OpenId   string
}

type TokenResp struct {
	AuthToken string
	OpenId    string
}

func GetUserInfo(openId, authToken string) (UserInfoResp, common.HttpErr) {
	var ret UserInfoResp
	result, er := tool2.HttpGet("http://localhost:4004/resource/api/v1/user/info", map[string]string{
		"open_id":      openId,
		"auth_token": authToken,
	})
	if er != nil {
		return ret, common.NewHttpErrWithError(er)
	}
	if result.Code != 0 {
		return ret, common.NewHttpErr(200, result.Code, result.Msg)
	}
	resp := result.Data.(map[string]interface{})
	ret.OpenId = resp["open_id"].(string)
	ret.Nickname = resp["nickname"].(string)
	ret.Face = resp["face"].(string)
	return ret, nil
}

func Certificate(code string) (TokenResp, common.HttpErr) {
	token := TokenResp{}
	result, er := tool2.HttpGet("http://localhost:4004/oauth/api/v1/certificate", map[string]string{
		"code": code,
	})
	if er != nil {
		return token, common.NewHttpErrWithError(er)
	}
	if result.Code != 0 {
		return token, common.NewHttpErr(200, result.Code, result.Msg)
	}
	resp := result.Data.(map[string]interface{})

	token.AuthToken = resp["auth_token"].(string)
	token.OpenId = resp["open_id"].(string)
	return token, nil
}
