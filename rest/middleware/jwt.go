package middleware

import (
	"github.com/cocoup/go-smart/rest/common/result"
	"github.com/cocoup/go-smart/rest/errorx"
	"github.com/cocoup/go-smart/rest/token"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/net/context"
	"net/http"
)

const (
	jwtAudience    = "aud"
	jwtExpire      = "exp"
	jwtId          = "jti"
	jwtIssueAt     = "iat"
	jwtIssuer      = "iss"
	jwtNotBefore   = "nbf"
	jwtSubject     = "sub"
	noDetailReason = "no detail reason"
)

func JWTAuth(secret string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tok, err := token.ParseToken(ctx.Request, secret)
		if err != nil {
			unauthorized(ctx, err)
			return
		}
		if !tok.Valid {
			result.HttpFailed(ctx, errorx.NewErrCode(errorx.TOKEN_EXPIRE))
			return
		}

		claims, ok := tok.Claims.(jwt.MapClaims)
		if !ok {
			result.HttpFailed(ctx, errorx.NewErrCode(errorx.TOKEN_CLAIMS_ERROR))
			return
		}

		reqCtx := ctx.Request.Context()
		for k, v := range claims {
			switch k {
			case jwtAudience, jwtExpire, jwtId, jwtIssueAt, jwtIssuer, jwtNotBefore, jwtSubject:
				// ignore the standard claims
			default:
				reqCtx = context.WithValue(ctx, k, v)
			}
		}

		ctx.Request = ctx.Request.WithContext(reqCtx)
		ctx.Set(token.KEY_TOKEN, tok)

		ctx.Next()
	}
}

func unauthorized(ctx *gin.Context, err error) {
	ctx.AbortWithStatus(http.StatusUnauthorized)
}
