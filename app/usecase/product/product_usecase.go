package product

import (
	"context"
	"orange-backstage-api/app/service/crawler"
	"orange-backstage-api/app/store"
	"orange-backstage-api/app/store/product"
)

type Usecase struct {
	store *store.Store
}

func New(store *store.Store) *Usecase {
	return &Usecase{store: store}
}

func (u *Usecase) UpsertProduct(ctx context.Context, newProduct crawler.Product) error {
	p := &product.Product{
		Name:  newProduct.Name,
		Types: product.Laptop,
	}

	pprice := &product.Price{
		Price:      newProduct.Price,
		SellerType: crawler.SellerTypeMOMO.String(),
		SourceURL:  newProduct.SourceURL,
	}

	return u.store.Product.UpsertProduct(ctx, p, pprice)
}
