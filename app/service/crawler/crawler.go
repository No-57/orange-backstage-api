package crawler

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
)

const randomImgURL = "https://picsum.photos/200/300"

type Crawler interface {
	// Fetch fetches products from the source.
	// TODO: More details about the source and conditions.
	Fetch(context.Context) ([]Product, error)

	SellerType() SellerType
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
	SellerTypeMOMO
)

func (st SellerType) String() string {
	switch st {
	case SellerTypePCHome:
		return "PCHOME"
	case SellerTypeShopee:
		return "SHOPEE"
	case SellerTypeAmazon:
		return "AMAZON"
	case SellerTypeYahoo:
		return "YAHOO"
	case SellerTypeOther:
		return "OTHER"
	case SellerTypeMOMO:
		return "MOMO"
	default:
		return "Unknown"
	}
}

type Product struct {
	ProductPrice
	ProductImg

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

type ProductImg struct {
	URL string
}
