package result

import (
	"fmt"
	"github.com/cocoup/go-smart/rest/errorx"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
)

func HttpSuccess(ctx *gin.Context, data interface{}) {
	ctx.Set("code", errorx.SUCCESS.String())
	ctx.JSON(http.StatusOK, Success(data))
}

func HttpFailed(ctx *gin.Context, err error) {
	//错误返回
	code := errorx.FAILED
	msg := "服务器开小差啦，请稍后再来试一试"

	causeErr := errors.Cause(err)              // err类型
	if e, ok := causeErr.(*errorx.Error); ok { //自定义错误类型
		code = e.GetCode()
		msg = e.GetMsg()
	}

	ctx.Set("code", code.String())
	ctx.JSON(http.StatusOK, Failed(uint(code), msg))
}

func HttpResult(ctx *gin.Context, data interface{}, err error) {
	if nil == err {
		ctx.Set("code", errorx.SUCCESS.String())
		if sData, ok := data.(string); ok {
			ctx.String(http.StatusOK, sData)
		} else {
			ctx.JSON(http.StatusOK, Success(data))
		}
	} else {
		HttpFailed(ctx, err)
	}
}

func ParamError(ctx *gin.Context, err error) {
	msg := fmt.Sprintf("%s ,%s", errorx.ErrMsg(errorx.PARAM_ERROR), err.Error())
	HttpFailed(ctx, errorx.NewErrCodeMsg(errorx.PARAM_ERROR, msg))
}
