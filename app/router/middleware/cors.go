package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CORS() gin.HandlerFunc {
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{
		"Origin",
		"Content-Length",
		"Content-Type",
		"Authorization",
		"Content-Type",
		"AccessToken",
		"X-CSRF-Token",
		"Authorization",
		"Token",
		"X-Token",
		"X-User-Id",
	}
	return cors.New(config)
}
