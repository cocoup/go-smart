package middleware

import (
	"bytes"
	"context"
	"fmt"
	"github.com/cocoup/go-smart/rest/log"
	"github.com/zeromicro/go-zero/core/color"
	"github.com/zeromicro/go-zero/rest/httpx"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zeromicro/go-zero/core/logx"

	"github.com/cocoup/go-smart/core/syncx"
	"github.com/cocoup/go-smart/core/timex"
	"github.com/cocoup/go-smart/core/utils"
	"github.com/cocoup/go-smart/rest/token"
)

const (
	limitBodyBytes       = 1024
	defaultSlowThreshold = time.Millisecond * 500
)

var slowThreshold = syncx.ForAtomicDuration(defaultSlowThreshold)

var bufferPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

type logResponseWriter struct {
	gin.ResponseWriter
	buf *bytes.Buffer
}

func (w logResponseWriter) Write(b []byte) (int, error) {
	w.buf.Write(b)
	return w.ResponseWriter.Write(b)
}

func LogHandler(verbose bool) gin.HandlerFunc {
	if verbose {
		return detailHandler()
	}
	return briefHandler()
}

// 不包含请求和响应体
func briefHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		timer := utils.NewElapsedTimer()
		logs := new(log.LogCollector)
		ctx.Request = ctx.Request.WithContext(context.WithValue(ctx.Request.Context(), log.LogContext, logs))
		ctx.Next()
		logBrief(ctx, timer, logs)
	}
}

func logBrief(ctx *gin.Context, timer *utils.ElapsedTimer, logs *log.LogCollector) {
	var buf bytes.Buffer
	duration := timer.Duration()
	httpCode := ctx.Writer.Status()
	logger := logx.WithContext(ctx.Request.Context()).WithDuration(duration)
	buf.WriteString(fmt.Sprintf("[HTTP] %s - %s %s - %s - %s",
		wrapStatusCode(httpCode), wrapMethod(ctx.Request.Method), ctx.Request.RequestURI, httpx.GetRemoteAddr(ctx.Request), ctx.Request.UserAgent()))
	if duration > slowThreshold.Load() {
		logger.Slowf("[HTTP] %s - %s - %s %s - %s - slowcall(%s)",
			wrapStatusCode(httpCode), wrapMethod(ctx.Request.Method), ctx.Request.RequestURI, httpx.GetRemoteAddr(ctx.Request), ctx.Request.UserAgent(),
			fmt.Sprintf("slowcall(%s)", timex.ReprOfDuration(duration)))
	}

	ok := httpCode < http.StatusInternalServerError
	if !ok {
		fullReq := dumpRequest(ctx.Request)
		limitReader := io.LimitReader(strings.NewReader(fullReq), limitBodyBytes)
		body, err := ioutil.ReadAll(limitReader)
		if err != nil {
			buf.WriteString(fmt.Sprintf("\n%s", fullReq))
		} else {
			buf.WriteString(fmt.Sprintf("\n%s", string(body)))
		}
	}

	body := logs.Flush()
	if len(body) > 0 {
		buf.WriteString(fmt.Sprintf("\n%s", body))
	}

	if ok {
		logger.Info(buf.String())
	} else {
		logger.Error(buf.String())
	}
}

func detailHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		timer := utils.NewElapsedTimer()

		var reqStr string
		reqBuff, err := ioutil.ReadAll(ctx.Request.Body)
		if nil != err {
			reqStr = err.Error()
		}
		ctx.Request.Body = ioutil.NopCloser(bytes.NewReader(reqBuff))
		reqStr = string(reqBuff)

		lrw := &logResponseWriter{buf: bufferPool.Get().(*bytes.Buffer), ResponseWriter: ctx.Writer}
		ctx.Writer = lrw

		logs := new(log.LogCollector)
		// 调用过程中收集信息(ctx.Value(log.LogContext).(*log.LogCollector).Append("anything"))
		ctx.Request = ctx.Request.WithContext(context.WithValue(ctx.Request.Context(), log.LogContext, logs))

		ctx.Next()

		logDetail(ctx, reqStr, lrw, timer, logs)
	}
}

func logDetail(ctx *gin.Context, reqStr string, lrw *logResponseWriter, timer *utils.ElapsedTimer, logs *log.LogCollector) {
	var buf bytes.Buffer

	duration := timer.Duration()
	httpCode := ctx.Writer.Status()
	path := ctx.Request.URL.Path
	remoteAddr := ctx.Request.RemoteAddr
	userAgent := ctx.Request.UserAgent()

	buf.WriteString(fmt.Sprintf("[HTTP] %s - %s - %d - %s - %s - %s\n=> %s\n",
		ctx.Request.Method, path, httpCode, remoteAddr, userAgent, timex.ReprOfDuration(duration), reqStr))

	body := logs.Flush()
	if len(body) > 0 {
		buf.WriteString(fmt.Sprintf("%s\n", body))
	}

	respBuf := lrw.buf.Bytes()
	if len(respBuf) > 0 {
		buf.WriteString(fmt.Sprintf("<= %s", respBuf))
	}
	lrw.buf.Reset()
	bufferPool.Put(lrw.buf)

	token, exists := ctx.Get(token.KEY_TOKEN)
	if exists {
		buf.WriteString(fmt.Sprintf("\n[token]%s", token.(string)))
	}

	logger := logx.WithContext(ctx.Request.Context())
	if httpCode < http.StatusInternalServerError {
		logger.Info(buf.String())
	} else {
		logger.Error(buf.String())
	}
}

func wrapMethod(method string) string {
	var colour color.Color
	switch method {
	case http.MethodGet:
		colour = color.BgBlue
	case http.MethodPost:
		colour = color.BgCyan
	case http.MethodPut:
		colour = color.BgYellow
	case http.MethodDelete:
		colour = color.BgRed
	case http.MethodPatch:
		colour = color.BgGreen
	case http.MethodHead:
		colour = color.BgMagenta
	case http.MethodOptions:
		colour = color.BgWhite
	}

	if colour == color.NoColor {
		return method
	}

	return logx.WithColorPadding(method, colour)
}

func wrapStatusCode(code int) string {
	var colour color.Color
	switch {
	case code >= http.StatusOK && code < http.StatusMultipleChoices:
		colour = color.BgGreen
	case code >= http.StatusMultipleChoices && code < http.StatusBadRequest:
		colour = color.BgBlue
	case code >= http.StatusBadRequest && code < http.StatusInternalServerError:
		colour = color.BgMagenta
	default:
		colour = color.BgYellow
	}

	return logx.WithColorPadding(strconv.Itoa(code), colour)
}

func dumpRequest(r *http.Request) string {
	reqContent, err := httputil.DumpRequest(r, true)
	if err != nil {
		return err.Error()
	}

	return string(reqContent)
}
