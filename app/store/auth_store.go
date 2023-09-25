package store

import (
	"context"
	"orange-backstage-api/app/store/auth"

	"gorm.io/gorm"
)

type Auth struct {
	db *gorm.DB
}

func (a Auth) CreateToken(ctx context.Context, token *auth.Token) error {
	return a.db.WithContext(ctx).Create(token).Error
}

func (a Auth) GetValidToken(
	ctx context.Context,
	accID uint64,
	aToken, rToken string,
) (*auth.Token, error) {
	var record *auth.Token
	if err := a.db.WithContext(ctx).
		Where(auth.Token{
			AccountID:    accID,
			AccessToken:  aToken,
			RefreshToken: rToken,
		}).
		Take(&record).Error; err != nil {
		return nil, err
	}
	return record, nil
}

func (a Auth) UpdateToken(ctx context.Context, token *auth.Token) error {
	return a.db.WithContext(ctx).Save(token).Error
}
