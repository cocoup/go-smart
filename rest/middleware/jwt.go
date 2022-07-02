package middleware

import (
	"github.com/cocoup/go-smart/core/jwt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func JWTAuth(conf jwt.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		/**
		我们这里jwt鉴权取头部信息 x-token 登录时回返回token信息
		这里前端需要把token存储到cookie或者本地localStorage中
		不过需要跟后端协商过期时间 可以约定刷新令牌或者重新登录
		*/
		token := ctx.Request.Header.Get("x-token")
		if token == "" {
			unauthorized(ctx, jwt.TokenInvalid)
			return
		}
		j := jwt.NewJWT(conf)
		// parseToken 解析token包含的信息
		claims, err := j.ParseToken(token)
		if err != nil {
			unauthorized(ctx, err)
			return
		}

		// 踢出逻辑
		//if service.Group.Sys.IsTokenValid(claims.UUID.String(), token) {
		//	common.RespFailedData(ctx, gin.H{"reload": true}, "您的帐户异地登陆或令牌失效")
		//	ctx.Abort()
		//	return
		//}

		// 用户被删除的逻辑 需要优化 此处比较消耗性能 如果需要 请自行打开
		//if err, _ = userService.FindUserByUuid(claims.UUID.String()); err != nil {
		//	_ = jwtService.JsonInBlacklist(system.JwtBlacklist{Jwt: token})
		//	response.FailWithDetailed(gin.H{"reload": true}, err.Error(), c)
		//	c.Abort()
		//}

		//if conf.Get().UseMultipoint {
		//	if claims.ExpiresAt.Unix()-time.Now().Unix() < claims.BufferTime {
		//		claims.ExpiresAt.Time = time.Unix(time.Now().Unix()+conf.ExpiresTime, 0)
		//		newToken, _ := j.CreateTokenByOldToken(token, *claims)
		//		newClaims, _ := j.ParseToken(newToken)
		//		ctx.Header("new-token", newToken)
		//		ctx.Header("new-expires-at", strconv.FormatInt(newClaims.ExpiresAt.Unix(), 10))
		//
		//		// 更新token
		//		_ = service.Group.Sys.ChangeToken(newClaims.UUID.String(), newToken)
		//	}
		//}
		ctx.Set("claims", claims)
		ctx.Next()
	}
}

func unauthorized(ctx *gin.Context, err error) {
	ctx.AbortWithStatus(http.StatusUnauthorized)
}
