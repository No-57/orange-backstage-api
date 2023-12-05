package router

import (
	"context"
	"orange-backstage-api/app/router/account"
	"orange-backstage-api/app/router/auth"
	"orange-backstage-api/app/router/board"
	"orange-backstage-api/app/router/image"
	"orange-backstage-api/app/router/middleware"
	"orange-backstage-api/app/router/theme"
	"orange-backstage-api/app/usecase"
	"orange-backstage-api/infra/api"
	"orange-backstage-api/infra/config"
	"time"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Router struct {
	ctx       context.Context
	version   string
	enableDoc bool
	param     Param

	auth    *auth.Router
	account *account.Router
	board   *board.Router
	image   *image.Router
	theme   *theme.Router
}

type Param struct {
	Version         string
	JWT             config.JWT
	EnableDoc       bool
	ImageUploadPath string
}

func New(ctx context.Context, usecase *usecase.Usecase, param Param) *Router {
	return &Router{
		ctx:     ctx,
		version: param.Version,
		param:   param,

		auth: auth.New(usecase.Auth),
		account: account.New(usecase.Account, account.Config{
			JWT: param.JWT,
		}),
		board: board.New(usecase.Board),
		image: image.New(image.Param{
			UploadPath: param.ImageUploadPath,
		}),
		theme: theme.New(usecase.Theme),

		enableDoc: param.EnableDoc,
	}
}

func (r *Router) Register(ginR gin.IRouter) {
	if r.enableDoc {
		r.registerSwagger(ginR)
	}

	api := ginR.Group("/api", generalMiddlewares(r.ctx)...)
	ver := api.Group("/" + r.version)

	r.image.Register(ver)
	ver.GET("/health", health)

	r.auth.Register(ver)
	r.account.Register(ver)

	auth := ver.Group("", middleware.JWTChceker(r.param.JWT.SecretBytes()))
	r.board.Register(auth)
	r.theme.Register(auth)
}

func (r *Router) registerSwagger(ginR gin.IRouter) {
	ginR.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

// health
//
//	@Summary		Health Check
//	@Description	Health Check
//	@Tags			health
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	router.health.resp
//	@Failure		400	{object}	api.ErrResp
//	@Failure		500	{object}	api.ErrResp
//	@Router			/health [get]
func health(c *gin.Context) {
	type resp struct {
		Health string    `json:"health" example:"ok"`
		Time   time.Time `json:"time" example:"2021-01-01T00:00:00+08:00"`
	}

	ctx := api.NewCtx(c)

	ctx.Resp().Data(resp{
		Health: "ok",
		Time:   time.Now().Local(),
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
