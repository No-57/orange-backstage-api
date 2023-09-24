package auth

import (
	"orange-backstage-api/app/usecase/auth"
	"orange-backstage-api/infra/api"

	"github.com/gin-gonic/gin"
)

type Router struct {
	usecase *auth.Usecase
}

func New(usecase *auth.Usecase) *Router {
	return &Router{
		usecase: usecase,
	}
}

func (r *Router) Register(ginR gin.IRouter) {
	ginR.POST("/login", r.Login)
}

type LoginPayload struct {
	Target   string `json:"target" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResp struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (r Router) Login(c *gin.Context) {
	api := api.New(c)

	var payload LoginPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		api.Resp().InvalidParam(err)
		return
	}

	token, err := r.usecase.Login(api.Ctx(), auth.LoginParam{
		Target:   payload.Target,
		Password: payload.Password,
	})
	if err != nil {
		api.Resp().Err(err)
		return
	}

	api.Resp().Data(LoginResp{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	})
}
