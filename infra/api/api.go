package api

import (
	"context"
	"orange-backstage-api/app/store/account"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

const APICtxKey = "api"

type API struct {
	c *gin.Context

	log  *zerolog.Logger
	resp *resp
	auth auth
}

func New(c *gin.Context) *API {
	if value, ok := c.Get(APICtxKey); ok {
		if api, ok := value.(*API); ok {
			return api
		}
	}

	log := zerolog.Ctx(c).With().Str("service", "api").Logger()
	api := &API{
		c:   c,
		log: &log,

		resp: newResp(c, log),
		auth: newAuth(c),
	}
	defer c.Set(APICtxKey, api)

	return api
}

func (api *API) Ctx() context.Context {
	return api.c.Request.Context()
}

func (api *API) Resp() *resp {
	return api.resp
}

func (api *API) Auth() auth {
	return api.auth
}

func (api *API) UpdateAuth(claims *account.Claims) {
	defer api.c.Set(APICtxKey, api)

	api.auth.Update(claims)
}
