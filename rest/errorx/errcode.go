package errorx

import "fmt"

type ErrCode uint

const (
	SUCCESS ErrCode = 0 //成功返回
	/**(前2位代表业务,后2位代表具体功能)**/
	FAILED               ErrCode = 10001 //失败返回
	PARAM_ERROR          ErrCode = 10002
	TOKEN_EXPIRE         ErrCode = 20001
	TOKEN_GENERATE_ERROR ErrCode = 20002
	DB_ERROR             ErrCode = 30001
)

func (e ErrCode) String() string {
	return fmt.Sprintf("%d", e)
}
