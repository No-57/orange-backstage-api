package auth

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"orange-backstage-api/app/store"
	"orange-backstage-api/app/store/auth"
	"orange-backstage-api/infra/config"
	"orange-backstage-api/infra/util/convert"

	"github.com/google/uuid"
)

type Usecase struct {
	store *store.Store
	param Param
}

type Param struct {
	JWT config.JWT
}

func New(store *store.Store, param Param) *Usecase {
	return &Usecase{
		store: store,
		param: param,
	}
}

type LoginParam struct {
	Target   string
	Password string
}

type LoginResponse struct {
	AccessToken  string
	RefreshToken string
}

func (u *Usecase) Login(
	ctx context.Context, param LoginParam,
) (*LoginResponse, error) {
	account, err := u.store.Account.GetByEmailOrName(ctx, param.Target)
	if err != nil {
		return nil, err
	}

	if ok := account.ComparePassword(param.Password); !ok {
		return nil, errors.New("invalid password")
	}

	jwtCfg := u.param.JWT
	accessToken, err := account.GenJWT(
		convert.StrToBytes(jwtCfg.Secret), jwtCfg.AccessTokenExpire,
	)
	if err != nil {
		return nil, fmt.Errorf("gen access token: %w", err)
	}
	refreshToken := geneRefreshToken()

	if err := u.store.Auth.CreateToken(ctx, &auth.Token{
		AccountID:    account.ID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}); err != nil {
		return nil, fmt.Errorf("create token: %w", err)
	}

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func geneRefreshToken() string {
	u := uuid.New()
	return fmt.Sprintf("%x", sha256.Sum256(u[:]))
}
