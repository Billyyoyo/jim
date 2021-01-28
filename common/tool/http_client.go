package tool

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"jim/common"
	"net/http"
)

func HttpGet(url string, params map[string]string) (result common.Result, err error) {
	sep := "?"
	var bt bytes.Buffer
	for k, v := range params {
		bt.Write([]byte(sep))
		bt.Write([]byte(k))
		bt.Write([]byte("="))
		bt.Write([]byte(v))
		sep = "&"
	}
	url = fmt.Sprintf("%s%s", url, bt.String())
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	if resp.StatusCode != http.StatusOK {
		err = errors.New(fmt.Sprintf("http call error_%d", resp.StatusCode))
		return
	}
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	result = common.Result{}
	err = json.Unmarshal(bs, &result)
	if err != nil {
		return
	}
	return
}
