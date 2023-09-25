package router

import (
	"context"
	"orange-backstage-api/app/router/account"
	"orange-backstage-api/app/router/auth"
	"orange-backstage-api/app/router/middleware"
	"orange-backstage-api/app/usecase"
	"orange-backstage-api/infra/api"
	"orange-backstage-api/infra/config"
	"time"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

type Router struct {
	ctx     context.Context
	version string

	auth    *auth.Router
	account *account.Router
}

type Param struct {
	Version string
	JWT     config.JWT
}

func New(ctx context.Context, usecase *usecase.Usecase, param Param) *Router {
	return &Router{
		ctx:     ctx,
		version: param.Version,

		auth: auth.New(usecase.Auth),
		account: account.New(usecase.Account, account.Config{
			JWT: param.JWT,
		}),
	}
}

func (r *Router) Register(ginR gin.IRouter) {
	api := ginR.Group("/api", generalMiddlewares(r.ctx)...)
	ver := api.Group("/" + r.version)

	ver.GET("/health", health)

	r.auth.Register(ver)
	r.account.Register(ver)
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
