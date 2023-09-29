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
	ginR.PATCH("/token", r.UpdateToken)
}

type LoginPayload struct {
	Target   string `json:"target" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResp struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// Login
//
//	@Summary		Account Login
//	@Description	Login with email or name and password
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			data	body		LoginPayload	true	"payload"
//	@Success		200		{object}	api.DataResp{data=LoginResp}
//	@Failure		400		{object}	api.ErrResp
//	@Failure		500		{object}	api.ErrResp
//	@Router			/login [post]
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

type UpdateTokenPayload struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// UpdateToken
//
//	@Summary		Update Token
//	@Description	Update Token with refresh token
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			data	body		UpdateTokenPayload	true	"payload"
//	@Success		200		{object}	api.DataResp{data=LoginResp}
//	@Failure		400		{object}	api.ErrResp
//	@Failure		500		{object}	api.ErrResp
//	@Router			/token [patch]
func (r Router) UpdateToken(c *gin.Context) {
	api := api.New(c)

	var payload UpdateTokenPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		api.Resp().InvalidParam(err)
		return
	}

	token, err := r.usecase.UpdateToken(api.Ctx(), auth.UpdateTokenParam{
		RToken: payload.RefreshToken,
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
