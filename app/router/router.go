package router

import (
	"context"
	"orange-backstage-api/app/router/auth"
	"orange-backstage-api/app/router/middleware"
	"orange-backstage-api/app/usecase"
	"orange-backstage-api/infra/api"
	"time"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

type Router struct {
	ctx     context.Context
	version string

	Auth *auth.Router
}

type Param struct {
	Version string
}

func New(ctx context.Context, usecase *usecase.Usecase, param Param) *Router {
	return &Router{
		ctx:     ctx,
		version: param.Version,

		Auth: auth.New(usecase.Auth),
	}
}

func (r *Router) Register(ginR gin.IRouter) {
	api := ginR.Group("/api", generalMiddlewares(r.ctx)...)
	ver := api.Group("/" + r.version)

	ver.GET("/health", health)

	r.Auth.Register(ver)
}

func health(c *gin.Context) {
	api := api.New(c)

	api.Resp().Data(gin.H{
		"health": "ok",
		"time":   time.Now().Local(),
	})
}

func generalMiddlewares(ctx context.Context) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		gin.Recovery(),
		middleware.CORS(),
		middleware.ReqID(),
		middleware.Logger(ctx),
		gzip.Gzip(gzip.DefaultCompression),
	}
}
