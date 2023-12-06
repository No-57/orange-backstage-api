package theme

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

func (s Store) Select(ctx context.Context) ([]model.Theme, error) {
	var records []model.Theme
	result := s.db.WithContext(ctx).Find(&records)
	return records, result.Error
}

func (s Store) SeleteByID(ctx context.Context, id uint64) (*model.Theme, error) {
	var record model.Theme
	result := s.db.WithContext(ctx).Take(&record, id)
	return &record, result.Error
}

func (s Store) Create(ctx context.Context, theme *model.Theme) error {
	result := s.db.WithContext(ctx).Create(theme)
	return result.Error
}

func (s Store) Delete(ctx context.Context, id ...uint64) error {
	result := s.db.WithContext(ctx).
		Where("id in ?", id).
		Delete(&model.Theme{})
	return result.Error
}
