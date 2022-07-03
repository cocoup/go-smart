package token

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"github.com/golang-jwt/jwt/v4/request"
	"net/http"
	"time"
)

var (
	TokenExpired     = errors.New("Token is expired")
	TokenNotValidYet = errors.New("Token not active yet")
	TokenMalformed   = errors.New("That's not even a token")
	TokenInvalid     = errors.New("Couldn't handle this token")
)

var KEY_TOKEN = "x-token"

func CreateToken(secret string, seconds int64, payloads map[string]interface{}) (string, error) {
	now := time.Now().Unix()
	claims := jwt.MapClaims{
		"exp": now + seconds,
		"iat": now,
	}
	for k, v := range payloads {
		claims[k] = v
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ParseToken(req *http.Request, secret string) (*jwt.Token, error) {
	return request.ParseFromRequest(req, request.AuthorizationHeaderExtractor,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		}, request.WithParser(jwt.NewParser(jwt.WithJSONNumber())))
}
