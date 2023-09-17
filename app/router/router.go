package router

import (
	"time"

	"github.com/gin-gonic/gin"
)

type Router struct {
	version string
}

type Param struct {
	Version string
}

func New(param Param) *Router {
	return &Router{
		version: param.Version,
	}
}

func (r *Router) Register(ginR gin.IRouter) {
	api := ginR.Group("/api")
	ver := api.Group("/" + r.version)

	ver.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"time":    time.Now().Local(),
			"message": "ok",
		})
	})
}
