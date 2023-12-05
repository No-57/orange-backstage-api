package theme

import (
	"context"
	"orange-backstage-api/app/model"
	"orange-backstage-api/app/store"
	"orange-backstage-api/infra/api"
)

type Usecase struct {
	store *store.Store
}

func New(store *store.Store) *Usecase {
	return &Usecase{store: store}
}

func (u *Usecase) List(ctx context.Context) ([]model.Theme, error) {
	records, err := u.store.Theme.Select(ctx)
	if err != nil {
		return nil, api.NewStoreErr(err)
	}

	return records, err
}

type CreateParam struct {
	Code    string
	Type    string
	Disable bool
}

func (u *Usecase) Create(ctx context.Context, param CreateParam) error {
	theme := &model.Theme{
		Code:    param.Code,
		Type:    param.Type,
		Disable: param.Disable,
	}
	if err := u.store.Theme.Create(ctx, theme); err != nil {
		return api.NewStoreErr(err)
	}

	return nil
}

func (u *Usecase) Delete(ctx context.Context, id uint64) error {
	if err := u.store.Theme.Delete(ctx, id); err != nil {
		return api.NewStoreErr(err)
	}

	return nil
}
