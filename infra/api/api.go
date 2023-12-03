package api

import (
	"context"
	"net/http"
	"orange-backstage-api/app/model"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

const APICtxKey = "api"

type Context struct {
	context.Context
	c *gin.Context

	log   *zerolog.Logger
	param *param
	resp  *resp
	auth  auth
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

		param: newParam(c),
		resp:  newResp(c, log),
		auth:  newAuth(c),
	}
	defer c.Set(APICtxKey, ctx)

	return ctx
}

func (ctx *Context) Param() *param {
	return ctx.param
}

func (ctx *Context) Resp() *resp {
	return ctx.resp
}

func (ctx *Context) Auth() auth {
	return ctx.auth
}

func (ctx *Context) UpdateAuth(claims *model.Claims) {
	defer ctx.c.Set(APICtxKey, ctx)

	ctx.auth.Update(claims)
}

func (ctx *Context) ShouldBind(obj interface{}) error {
	if err := ctx.c.ShouldBind(obj); err != nil {
		return NewParamErr(err)
	}

	return nil
}

func (ctx *Context) ShouldBindJSON(obj interface{}) error {
	if err := ctx.c.ShouldBindJSON(obj); err != nil {
		return NewParamErr(err)
	}

	return nil
}

func (ctx *Context) ShouldBindUri(obj interface{}) error {
	if err := ctx.c.ShouldBindUri(obj); err != nil {
		return NewHTTPErr(http.StatusNotFound, CodeAPINotFound, err)
	}

	return nil
}
