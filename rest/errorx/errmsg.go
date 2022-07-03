package errorx

var errMsg map[ErrCode]string

func init() {
	errMsg = make(map[ErrCode]string)
	errMsg[SUCCESS] = "SUCCESS"
	errMsg[FAILED] = "服务器开小差啦,稍后再来试一试"
	errMsg[DATA_NOT_FOUND] = "数据未找到"
	errMsg[PARAM_ERROR] = "参数错误"
	errMsg[TOKEN_EXPIRE] = "token失效，请重新登陆"
	errMsg[TOKEN_GENERATE_ERROR] = "生成token失败"
	errMsg[TOKEN_PARSE_ERROR] = "解析token失败"
	errMsg[DB_ERROR] = "数据库繁忙,请稍后再试"
}

func ErrMsg(code ErrCode) string {
	if msg, ok := errMsg[code]; ok {
		return msg
	} else {
		return "服务器开小差啦,请稍后再来试一试"
	}
}

func IsErrCode(code ErrCode) bool {
	_, ok := errMsg[code]
	return ok
}
