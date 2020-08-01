package core

import (
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strconv"
)

const (
	API_URL = "http://localhost:4001/jim/api/v1"
)

type Result struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Desc string      `json:"desc"`
	Data interface{} `json:"data"`
}

type AuthInfo struct {
	TcpServer string `json:"tcpServer"`
	Uid       int64  `json:"uid"`
	Token     string `json:"token"`
}

type MUser struct {
	Uid  int64
	Name string
}

type MSession struct {
	Sid  int64
	Name string
	Type int8
}

func Authorization(code string) (uid int64, token string, server string, err error) {
	params := fmt.Sprintf("code=%s", code)
	resp, err := http.Get(API_URL + "/enter?" + params)
	if err != nil {
		return
	}
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	log.Info("call http /enter, response:\n", string(bs))
	result := &Result{}
	err = json.Unmarshal(bs, result)
	if err != nil {
		return
	}
	if result.Code > 0 {
		err = errors.New(fmt.Sprintf("[%d]:%s(%s)", result.Code, result.Msg, result.Desc))
		return
	}
	m := result.Data.(map[string]interface{})
	uid = int64(m["uid"].(float64))
	token = m["token"].(string)
	server = m["tcpServer"].(string)
	return
}

func setHeader(req *http.Request, ctx *Ctx) {
	req.Header.Add("jim-uid", strconv.FormatInt(ctx.UserId, 10))
	req.Header.Add("jim-device", strconv.FormatInt(ctx.DeviceId, 10))
	req.Header.Add("jim-token", ctx.Token)
}

func GetSessions(ctx *Ctx) (sessions *[]MSession, err error) {
	req, err := http.NewRequest("GET", API_URL+"/session/list", nil)
	if err != nil {
		log.Error("get sessions - new req fail: ", err.Error())
		return
	}
	setHeader(req, ctx)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Error("get sessions - do fail: ", err.Error())
		return
	}
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	log.Info("call http /session/list, response:\n", string(bs))
	result := &Result{}
	err = json.Unmarshal(bs, result)
	if err != nil {
		log.Error(err.Error())
		return
	}
	sessions = &[]MSession{}
	ms := result.Data.([]interface{})
	for _, m := range ms {
		itf := m.(map[string]interface{})
		s := MSession{
			Sid:  int64(itf["id"].(float64)),
			Name: itf["name"].(string),
			Type: int8(itf["type"].(float64)),
		}
		*sessions = append(*sessions, s)
	}
	return
}

func GetMembers(ctx *Ctx) (members *[]MUser, err error) {
	req, err := http.NewRequest("GET", fmt.Sprintf(API_URL+"/session/members?sessionId=%d", ctx.SessionId), nil)
	if err != nil {
		log.Error("get members - new req fail: ", err.Error())
		return
	}
	setHeader(req, ctx)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Error("get sessions - do fail: ", err.Error())
		return
	}
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	log.Info("call http /session/list, response:\n", string(bs))
	result := &Result{}
	err = json.Unmarshal(bs, result)
	if err != nil {
		log.Error(err.Error())
		return
	}
	members = &[]MUser{}
	ms := result.Data.([]interface{})
	for _, m := range ms {
		mm := m.(map[string]interface{})
		member := MUser{
			Uid:  int64(mm["id"].(float64)),
			Name: mm["name"].(string),
		}
		*members = append(*members, member)
	}
	return
}

func GetSession(ctx *Ctx, sessionId int64) (session *MSession, members *[]MUser, err error) {
	req, err := http.NewRequest("GET", fmt.Sprintf(API_URL+"/session/get?sessionId=%d", ctx.SessionId), nil)
	if err != nil {
		log.Error("get session - new req fail: ", err.Error())
		return
	}
	setHeader(req, ctx)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Error("get session - do fail: ", err.Error())
		return
	}
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	log.Info("call http /session/get, response:\n", string(bs))
	result := &Result{}
	err = json.Unmarshal(bs, result)
	if err != nil {
		log.Error(err.Error())
		return
	}
	m := result.Data.(map[string]interface{})
	sm := m["session"].(map[string]interface{})
	session = &MSession{
		Sid:  int64(sm["id"].(float64)),
		Name: sm["name"].(string),
		Type: int8(sm["type"].(float64)),
	}
	ms := m["members"].([]interface{})
	members = &[]MUser{}
	for _, mem := range ms {
		mm := mem.(map[string]interface{})
		member := MUser{
			Uid:  int64(mm["id"].(float64)),
			Name: mm["name"].(string),
		}
		*members = append(*members, member)
	}
	return
}
