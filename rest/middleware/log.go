package middleware

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"sync"
)

var bufferPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	if n, err := w.body.Write(b); err != nil {
		return n, err
	}
	return w.ResponseWriter.Write(b)
}

func LogHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		//start := time.Now()
		//path := c.Request.URL.Path
		//if 0 <= strings.Index(path, `metrics`) {
		//	c.Next()
		//	return
		//}

		path := c.Request.URL.Path
		fmt.Println(path)

		c.Next()
	}
}
