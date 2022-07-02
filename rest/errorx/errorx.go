package errorx

import (
	"fmt"
)

type Error struct {
	code ErrCode
	msg  string
}

func (e *Error) Error() string {
	return fmt.Sprintf("ErrCode:%d，ErrMsg:%s", e.code, e.msg)
}

func (e *Error) GetCode() ErrCode {
	return e.code
}

func (e *Error) GetMsg() string {
	return e.msg
}

func NewErrCodeMsg(code ErrCode, msg string) *Error {
	return &Error{code: code, msg: msg}
}

func NewErrCode(code ErrCode) *Error {
	return &Error{code: code, msg: ErrMsg(code)}
}

func NewErrMsg(msg string) *Error {
	return &Error{code: FAILED, msg: msg}
}
