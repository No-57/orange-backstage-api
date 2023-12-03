package auth

import (
	"context"
	"orange-backstage-api/app/model"

	"gorm.io/gorm"
)

type Store struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Store {
	return &Store{db: db}
}

func (s Store) CreateToken(ctx context.Context, token *model.Token) error {
	return s.db.WithContext(ctx).Create(token).Error
}

func (s Store) GetValidToken(
	ctx context.Context,
	accID uint64,
	aToken, rToken string,
) (*model.Token, error) {
	var record *model.Token
	if err := s.db.WithContext(ctx).
		Where(model.Token{
			AccountID:    accID,
			AccessToken:  aToken,
			RefreshToken: rToken,
		}).
		Take(&record).Error; err != nil {
		return nil, err
	}
	return record, nil
}

func (s Store) UpdateToken(ctx context.Context, token *model.Token) error {
	return s.db.WithContext(ctx).Save(token).Error
}
