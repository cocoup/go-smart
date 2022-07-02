package middleware

import (
	"fmt"
	"github.com/cocoup/go-smart/core/prometheus"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"time"
)

func Prometheus(monitor *prometheus.Monitor, engine *gin.Engine) gin.HandlerFunc {
	engine.GET("/metrics", gin.WrapH(promhttp.Handler()))

	return func(c *gin.Context) {
		//使用rest风格api路径时可结合group_wrapper使用
		//relativePath := c.GetString(constant.RelativePathKey)
		//if relativePath == "" {
		//	relativePath = c.Request.URL.Path
		//}

		relativePath := c.Request.URL.Path
		start := time.Now()
		reqSize := computeApproximateRequestSize(c.Request)
		c.Next()
		duration := time.Since(start)

		httpCode := fmt.Sprintf("%d", c.Writer.Status())
		code := "0"
		if val, exists := c.Get("code"); exists {
			code = val.(string)
		}
		errMsg := ""
		if val, exists := c.Get("err_msg"); exists {
			errMsg = val.(string)
		}

		monitor.With(
			prometheus.APIRequestsCounter(prometheus.Labels{"handler": relativePath, "method": c.Request.Method, "http_code": httpCode, "code": code, "err_msg": errMsg}),
			prometheus.RequestDuration(prometheus.Labels{"handler": relativePath, "method": c.Request.Method, "http_code": httpCode, "code": code}, duration.Seconds()),
			prometheus.RequestSize(prometheus.Labels{"handler": relativePath, "method": c.Request.Method, "http_code": httpCode, "code": code}, float64(reqSize)),
			prometheus.ResponseSize(prometheus.Labels{"handler": relativePath, "method": c.Request.Method, "http_code": httpCode, "code": code}, float64(c.Writer.Size())),
		)
	}
}

// From https://github.com/DanielHeckrath/gin-prometheus/blob/master/gin_prometheus.go
func computeApproximateRequestSize(r *http.Request) int {
	s := 0
	if r.URL != nil {
		s = len(r.URL.Path)
	}

	s += len(r.Method)
	s += len(r.Proto)
	for name, values := range r.Header {
		s += len(name)
		for _, value := range values {
			s += len(value)
		}
	}
	s += len(r.Host)

	// N.B. r.Form and r.MultipartForm are assumed to be included in r.URL.

	if r.ContentLength != -1 {
		s += int(r.ContentLength)
	}
	return s
}
