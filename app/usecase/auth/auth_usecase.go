package auth

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"orange-backstage-api/app/store"
	"orange-backstage-api/app/store/auth"
	"orange-backstage-api/infra/config"
	"time"

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
		jwtCfg.SecretBytes(), jwtCfg.AccessTokenExpire,
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

type UpdateTokenParam struct {
	AccID          uint64
	AToken, RToken string
}

func (u *Usecase) UpdateToken(
	ctx context.Context, param UpdateTokenParam,
) (*LoginResponse, error) {
	token, err := u.store.Auth.GetValidToken(
		ctx, param.AccID, param.AToken, param.RToken,
	)
	if err != nil {
		return nil, fmt.Errorf("get token: %w", err)
	}

	jwtCfg := u.param.JWT
	if time.Since(token.UpdatedAt) > jwtCfg.RefreshTokenExpire {
		return nil, errors.New("refresh token expired")
	}

	account, err := u.store.Account.GetByID(ctx, token.AccountID)
	if err != nil {
		return nil, fmt.Errorf("get account: %w", err)
	}

	accessToken, err := account.GenJWT(
		jwtCfg.SecretBytes(), jwtCfg.AccessTokenExpire,
	)
	if err != nil {
		return nil, fmt.Errorf("gen access token: %w", err)
	}

	token.AccessToken = accessToken
	token.RefreshToken = geneRefreshToken()
	if err := u.store.Auth.UpdateToken(ctx, token); err != nil {
		return nil, fmt.Errorf("update token: %w", err)
	}

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: token.RefreshToken,
	}, nil
}

func geneRefreshToken() string {
	u := uuid.New()
	return fmt.Sprintf("%x", sha256.Sum256(u[:]))
}
