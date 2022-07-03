package middleware

import (
	"compress/gzip"
	"github.com/gin-gonic/gin"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
	"strings"
)

func Unzip() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if strings.Contains(ctx.Request.Header.Get(httpx.ContentEncoding), "gzip") {
			reader, err := gzip.NewReader(ctx.Request.Body)
			if err != nil {
				ctx.AbortWithStatus(http.StatusBadRequest)
				return
			}

			ctx.Request.Body = reader
			ctx.Next()
		}
	}
}
