package crawler

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
)

type Crawler interface {
	// Fetch fetches products from the source.
	// TODO: More details about the source and conditions.
	Fetch(context.Context) ([]Product, error)
}

type ProductType int

const (
	ProductTypeUnknown ProductType = iota
	ProductTypePhone
	ProductTypeLaptop
	ProductTypeDesktop
	ProductTypeAudio
	ProductTypeTablet
	ProductTypeEarphone
)

type SellerType int

const (
	SellerTypeUnknown SellerType = iota
	SellerTypePCHome
	SellerTypeShopee
	SellerTypeAmazon
	SellerTypeYahoo
	SellerTypeOther
)

type Product struct {
	ProductPrice

	ID          uint64
	Name        string
	Description string
	Type        ProductType

	UpdatedAt time.Time
}

type ProductPrice struct {
	Price      decimal.Decimal
	Discount   decimal.Decimal
	SourceURL  string
	SellerType SellerType
}
