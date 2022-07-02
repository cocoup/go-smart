package middleware

import "github.com/gin-gonic/gin"

func {{.name}}() gin.HandlerFunc {
	return func(ctx *gin.Context) {
        // TODO generate middleware implement function, delete after code implementationv

        ctx.Next()
	}
}
