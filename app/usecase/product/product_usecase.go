package product

import (
	"context"
	"orange-backstage-api/app/model"
	"orange-backstage-api/app/service/crawler"
	"orange-backstage-api/app/store"
)

type Usecase struct {
	store *store.Store
}

func New(store *store.Store) *Usecase {
	return &Usecase{store: store}
}

func (u *Usecase) UpsertProduct(ctx context.Context, newProduct crawler.Product) error {
	p := &model.Product{
		Name:  newProduct.Name,
		Types: model.Laptop,
	}

	pprice := &model.Price{
		Price:      newProduct.Price,
		SellerType: crawler.SellerTypeMOMO.String(),
		Discount:   newProduct.Discount,
		SourceURL:  newProduct.SourceURL,
	}

	return u.store.Product.UpsertProduct(ctx, p, pprice)
}
