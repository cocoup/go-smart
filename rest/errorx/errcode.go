package errorx

import "fmt"

type ErrCode uint

const (
	SUCCESS ErrCode = 0 //成功返回
	/**(前2位代表业务,后2位代表具体功能)**/
	FAILED               ErrCode = 1001 //失败返回
	PARAM_ERROR          ErrCode = 1002
	DATA_NOT_FOUND       ErrCode = 1003 //数据未找到
	TOKEN_EXPIRE         ErrCode = 2001
	TOKEN_GENERATE_ERROR ErrCode = 2002
	TOKEN_CLAIMS_ERROR   ErrCode = 2003
	DB_ERROR             ErrCode = 3001
)

func (e ErrCode) String() string {
	return fmt.Sprintf("%d", e)
}
