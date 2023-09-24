package api

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type API struct {
	c *gin.Context

	log  *zerolog.Logger
	resp *Resp
}

func New(c *gin.Context) *API {
	log := zerolog.Ctx(c).With().Str("service", "api").Logger()

	return &API{
		c: c,

		resp: NewResp(c, log),
		log:  &log,
	}
}

func (api *API) Ctx() context.Context {
	return api.c.Request.Context()
}

func (api *API) Resp() *Resp {
	return api.resp
}
