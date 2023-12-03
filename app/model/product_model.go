package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type Category string

func (c Category) String() string {
	return string(c)
}

const (
	Phone   Category = "phone"
	Laptop  Category = "laptop"
	Desktop Category = "desktop"
	Audio   Category = "audio"
	Tablet  Category = "tablet"
	Watch   Category = "watch"
)

type Product struct {
	ID          uint64    `gorm:"column:id;primaryKey"`
	Name        string    `gorm:"column:name;NOT NULL"`
	Description string    `gorm:"column:description;NOT NULL"`
	Types       Category  `gorm:"column:types;NOT NULL"`
	CreatedAt   time.Time `gorm:"column:created_date;NOT NULL"`
	UpdatedAt   time.Time `gorm:"column:updated_date;NOT NULL"`
}

type Price struct {
	ID         uint64          `gorm:"column:id;primaryKey"`
	ProductID  uint64          `gorm:"column:product_id;NOT NULL"`
	Price      decimal.Decimal `gorm:"column:price;NOT NULL"`
	Discount   decimal.Decimal `gorm:"column:discount;NOT NULL"`
	SourceURL  string          `gorm:"column:source_url;NOT NULL"`
	SellerType string          `gorm:"column:seller_type;NOT NULL"`
	UpdatedAt  time.Time       `gorm:"column:updated_date;NOT NULL"`
}

func (p Price) TableName() string {
	return "price"
}

