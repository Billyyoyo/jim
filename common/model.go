package common

import "fmt"

type Result struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Desc string      `json:"desc"`
	Data interface{} `json:"data"`
}

type httpErr struct {
	status  int
	code    int
	message string
}

func (e httpErr) Status() int {
	return e.status
}
func (e httpErr) Code() int {
	return e.code
}
func (e httpErr) Message() string {
	return e.message
}

type HttpErr interface {
	error
	Status() int
	Code() int
	Message() string
}

func NewHttpErr(status, code int, msg ...string) HttpErr {
	err := httpErr{
		status: status,
		code:   code,
	}
	if msg != nil && len(msg) > 0 {
		err.message = msg[0]
	}
	return err
}

func NewHttpErrWithError(err error) HttpErr {
	return httpErr{
		message: err.Error(),
	}
}

func (e httpErr) Error() string {
	return fmt.Sprintf("status=%d\tcode=%d\tmessage=%s", e.status, e.code, e.message)
}
