package store

import (
	"context"
	"orange-backstage-api/app/store/product"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Product struct {
	db *gorm.DB
}

func (p *Product) UpsertProduct(ctx context.Context, product *product.Product, price *product.Price) error {
	result := p.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "name"}},
		UpdateAll: true,
	}).Create(product)
	if err := result.Error; err != nil {
		return err
	}

	result = p.db.Take(product)
	if err := result.Error; err != nil {
		return err
	}

	price.ProductID = product.ID
	result = p.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "product_id"}},
		UpdateAll: true,
	}).Create(price)

	return result.Error
}
