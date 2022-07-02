package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/zeromicro/go-zero/core/logx"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
)

func Recovery(stack bool) gin.HandlerFunc {
	return RecoveryWithWriter(stack)
}

// Recovery recover掉项目可能出现的panic
func RecoveryWithWriter(stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") ||
							strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					logx.Error(c.Request.URL.Path,
						"[error]", err,
						"[request]", string(httpRequest),
					)
					// If the connection is dead, we can't write a status to it.
					_ = c.Error(err.(error)) // nolint: errcheck
					c.Abort()
					return
				}

				if stack {
					logx.Error("[Recovery from panic]", err,
						"[request]", string(httpRequest),
						"[stack]", string(debug.Stack()),
					)
				} else {
					logx.Error("[Recovery from panic]",
						"[error]", err,
						"[request]", string(httpRequest),
					)
				}
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
