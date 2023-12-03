package product

import (
	"context"
	"orange-backstage-api/app/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Store struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Store {
	return &Store{db: db}
}

func (s Store) UpsertProduct(ctx context.Context, product *model.Product, price *model.Price) error {
	result := s.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "name"}},
		UpdateAll: true,
	}).Create(product)
	if err := result.Error; err != nil {
		return err
	}

	result = s.db.Take(product)
	if err := result.Error; err != nil {
		return err
	}

	price.ProductID = product.ID
	result = s.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "product_id"}},
		UpdateAll: true,
	}).Create(price)

	return result.Error
}
