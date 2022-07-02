package result

import "github.com/cocoup/go-smart/rest/errorx"

type Response struct {
	Code uint        `json:"code"`
	Data interface{} `json:"data,omitempty"`
	Msg  string      `json:"msg,omitempty"`
}

func Success(data interface{}) *Response {
	return &Response{Code: uint(errorx.SUCCESS), Data: data}
}

func Failed(errCode uint, errMsg string) *Response {
	return &Response{Code: errCode, Msg: errMsg}
}
