package api

import (
	"context"
	"orange-backstage-api/app/model"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

const APICtxKey = "api"

type Context struct {
	context.Context
	c *gin.Context

	log  *zerolog.Logger
	resp *resp
	auth auth
}

func NewCtx(c *gin.Context) *Context {
	if value, ok := c.Get(APICtxKey); ok {
		if api, ok := value.(*Context); ok {
			return api
		}
	}

	log := zerolog.Ctx(c).With().Str("service", "api").Logger()
	ctx := &Context{
		Context: c.Request.Context(),

		c:   c,
		log: &log,

		resp: newResp(c, log),
		auth: newAuth(c),
	}
	defer c.Set(APICtxKey, ctx)

	return ctx
}

func (api *Context) Resp() *resp {
	return api.resp
}

func (api *Context) Auth() auth {
	return api.auth
}

func (api *Context) UpdateAuth(claims *model.Claims) {
	defer api.c.Set(APICtxKey, api)

	api.auth.Update(claims)
}
