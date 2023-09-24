package usecase

import (
	"orange-backstage-api/app/store"
	"orange-backstage-api/app/usecase/auth"
	"orange-backstage-api/infra/config"
)

type Usecase struct {
	store *store.Store

	Auth *auth.Usecase
}

type Config struct {
	JWT config.JWT
}

func New(store *store.Store, config Config) *Usecase {
	return &Usecase{
		store: store,

		Auth: auth.New(store, auth.Param{
			JWT: config.JWT,
		}),
	}
}
