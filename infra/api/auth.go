package api

import (
	"orange-backstage-api/app/model"
	"strconv"

	"github.com/gin-gonic/gin"
)

const AuthCtxKey = "auth"

type auth struct {
	claims *model.Claims

	accID uint64
}

func newAuth(c *gin.Context) auth {
	var auth auth

	value, ok := c.Get(AuthCtxKey)
	if !ok {
		return auth
	}

	claims, ok := value.(*model.Claims)
	if !ok {
		return auth
	}
	auth.claims = claims

	accID, err := strconv.ParseUint(claims.Subject, 10, 64)
	if err != nil {
		return auth
	}
	auth.accID = accID

	return auth
}

func (a auth) Exists() bool {
	return a.claims != nil
}

func (a *auth) Update(claims *model.Claims) {
	a.claims = claims
}

func (a auth) AccID() uint64 {
	return a.accID
}
