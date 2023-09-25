package middleware

import (
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
)

func ReqID() gin.HandlerFunc {
	return requestid.New(requestid.WithGenerator(
		func() string {
			return xid.New().String()
		},
	))
}
