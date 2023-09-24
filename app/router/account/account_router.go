package account

import (
	"orange-backstage-api/app/router/middleware"
	"orange-backstage-api/app/usecase/account"
	"orange-backstage-api/infra/api"
	"orange-backstage-api/infra/config"
	"time"

	"github.com/gin-gonic/gin"
)

type Router struct {
	usecase *account.Usecase
	cfg     Config
}

type Config struct {
	JWT config.JWT
}

func New(usecase *account.Usecase, cfg Config) *Router {
	return &Router{
		usecase: usecase,
		cfg:     cfg,
	}
}

func (r *Router) Register(ginR gin.IRouter) {
	auth := ginR.Use(
		middleware.JWTChceker(r.cfg.JWT.SecretBytes()),
	)

	auth.GET("/self", r.Self)
}

type SelfResp struct {
	ID        uint64    `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (r Router) Self(c *gin.Context) {
	api := api.New(c)

	accID := api.Auth().AccID()
	acc, err := r.usecase.Self(api.Ctx(), account.SelfParam{ID: accID})
	if err != nil {
		api.Resp().Err(err)
		return
	}

	api.Resp().Data(SelfResp{
		ID:        acc.ID,
		Email:     acc.Email,
		Name:      acc.Name,
		CreatedAt: acc.CreatedAt,
		UpdatedAt: acc.UpdatedAt,
	})
}
