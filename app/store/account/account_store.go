package account

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

func (s Store) GetByEmailOrName(ctx context.Context, val string) (*model.Account, error) {
	var record *model.Account
	if err := s.db.WithContext(ctx).
		Where(model.Account{Email: val}).
		Or(model.Account{Name: val}).
		Take(&record).Error; err != nil {
		return nil, err
	}
	return record, nil
}

func (s Store) GetByID(ctx context.Context, ID uint64) (*model.Account, error) {
	var record *model.Account
	if err := s.db.WithContext(ctx).
		Where(model.Account{ID: ID}).
		Take(&record).Error; err != nil {
		return nil, err
	}
	return record, nil
}
