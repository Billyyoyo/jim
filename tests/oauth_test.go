package tests

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"jim/common/utils"
	"net/url"

	//tool2 "jim/http/tool"
	"jim/oauth/dao"
	"jim/oauth/model"
	"testing"
	"time"
)

func TestGetUserByLoginName(t *testing.T) {
	user, err := dao.GetUserByLoginName("billyyoyo")
	if err != nil {
		println(err.Error())
		return
	}
	printj(user)
}

func TestRandomString(t *testing.T) {
	for i := 0; i < 100; i++ {
		str := utils.RandomStr(8)
		fmt.Println(str)
	}
}

func TestVerifySigniture(t *testing.T) {
	word, err := utils.SignRsa(utils.PUBLIC_KEY, []byte("hello world"))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("Sign text:", word)
	current := time.Now().UnixNano()
	err = utils.VerifyRsa(utils.PRIVATE_KEY, []byte("hello world"), word)
	if err != nil {
		fmt.Println("verify failed")
	} else {
		fmt.Println("Success", (time.Now().UnixNano() - current))
	}
}

func TestDESCrypto(t *testing.T) {
	current := time.Now().UnixNano()
	word := "hello world|abc"
	bytes := utils.EncryptDES([]byte(word), []byte(utils.DES_KEY))
	fmt.Println(string(bytes))
	bytes = utils.DecryptDES(bytes, []byte(utils.DES_KEY))
	fmt.Println(string(bytes))
	word = "hello world|efg"
	bytes = utils.EncryptDES([]byte(word), []byte(utils.DES_KEY))
	fmt.Println(string(bytes))
	fmt.Println(time.Now().UnixNano() - current)
}

func TestRSACrypto(t *testing.T) {
	current := time.Now().UnixNano()
	word := "hello world|abc"
	bytes, _ := utils.EncryptRSA([]byte(word), utils.PUBLIC_KEY)
	fmt.Println(base64.StdEncoding.EncodeToString(bytes))
	fmt.Println(time.Now().UnixNano() - current)
	word = "hello world|def"
	bytes, _ = utils.EncryptRSA([]byte(word), utils.PUBLIC_KEY)
	fmt.Println(base64.StdEncoding.EncodeToString(bytes))
}

func TestHMacSign(t *testing.T) {
	word := "hello world|abc"
	bytes := utils.SignHmac([]byte(word), []byte(utils.HMAC_KEY))
	fmt.Println(base64.StdEncoding.EncodeToString(bytes))
	word = "hello world|efg"
	bytes = utils.SignHmac([]byte(word), []byte(utils.HMAC_KEY))
	fmt.Println(base64.StdEncoding.EncodeToString(bytes))
}

func TestAuthCode(t *testing.T) {
	curr := time.Now().UnixNano()
	authCode := model.AuthCode{
		OpenId:        "1001",
		Expired:       utils.GetCurrentMS() + 180000,
		ApplicationId: 6000,
	}
	bytes, err := json.Marshal(authCode)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	bytes = utils.EncryptDES(bytes, []byte(utils.DES_KEY))
	fmt.Println(base64.StdEncoding.EncodeToString(bytes))
	fmt.Println(time.Now().UnixNano() - curr)
	bytes = utils.SignHmac(bytes, []byte(utils.HMAC_KEY))
	fmt.Println(base64.StdEncoding.EncodeToString(bytes))
	fmt.Println(time.Now().UnixNano() - curr)
}

func TestArray(t *testing.T) {
	testArr()
}

func testArr(values ...string) {
	if values == nil {
		fmt.Println("nil")
	}
}

func TestUrlEncode(t *testing.T) {
	str := "http://localhost:4001/auth/callback"
	text := url.QueryEscape(str)
	println(text)
}
