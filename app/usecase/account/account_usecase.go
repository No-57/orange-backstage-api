package account

import (
	"context"
	"orange-backstage-api/app/store"
	"time"
)

type Usecase struct {
	store *store.Store
}

func New(store *store.Store) *Usecase {
	return &Usecase{
		store: store,
	}
}

type SelfParam struct {
	ID uint64
}

type SelfResponse struct {
	ID        uint64
	Email     string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (u *Usecase) Self(
	ctx context.Context, param SelfParam,
) (*SelfResponse, error) {
	account, err := u.store.Account.GetByID(ctx, param.ID)
	if err != nil {
		return nil, err
	}

	return &SelfResponse{
		ID:        account.ID,
		Email:     account.Email,
		Name:      account.Name,
		CreatedAt: account.CreatedAt,
		UpdatedAt: account.UpdatedAt,
	}, nil
}
