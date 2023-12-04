package api

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

type param struct {
	c *gin.Context
}

func newParam(c *gin.Context) *param {
	return &param{c}
}

func (p param) Get(key string) string {
	return p.c.Param(key)
}

func (p param) ID() uint64 {
	value := p.Get("id")
	id, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return 0
	}

	return id
}
