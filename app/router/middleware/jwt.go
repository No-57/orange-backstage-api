package middleware

import (
	"errors"
	"orange-backstage-api/app/model"
	"orange-backstage-api/infra/api"
	"strings"

	"github.com/gin-gonic/gin"
)

func JWTChceker(secret []byte) gin.HandlerFunc {
	const (
		jwtHeaderKey = "Authorization"
		bearer       = "Bearer "
		ctxKey       = api.AuthCtxKey
	)

	return func(c *gin.Context) {
		if value, ok := c.Get(api.AuthCtxKey); ok {
			if _, ok := value.(*model.Claims); ok {
				c.Next()
				return
			}
		}

		api := api.NewCtx(c)
		if api.Auth().Exists() {
			c.Next()
			return
		}

		header := c.GetHeader(jwtHeaderKey)
		value := strings.Split(header, bearer)
		if len(value) != 2 {
			api.Resp().Forbidden(errors.New("invalid header"))
			return
		}

		token := value[1]
		if token == "" {
			api.Resp().Forbidden(errors.New("invalid token"))
			return
		}

		claims, err := model.ParseTokenWithSecret(secret, token)
		if err != nil {
			if errors.Is(err, model.ErrExpired) {
				api.Resp().ExpiredToken(err)
				return
			}

			api.Resp().Forbidden(errors.New("invalid token"))
			return
		}

		api.UpdateAuth(claims)

		c.Next()
	}
}
