package service

import (
	"context"
	"orange-backstage-api/app/service/crawler"
	"orange-backstage-api/app/usecase"

	"github.com/rs/zerolog"
)

type Service struct {
	Crawlers []crawler.Crawler

	usecase *usecase.Usecase
}

func New(usecase *usecase.Usecase) *Service {
	crawlers := []crawler.Crawler{
		crawler.NewMomo(),
		crawler.NewPchome(),
	}
	return &Service{
		crawlers,
		usecase,
	}
}

func (s *Service) Run(ctx context.Context) error {
	// TODO: cron job
	for _, c := range s.Crawlers {
		products, err := c.Fetch(ctx)
		if err != nil {
			zerolog.Ctx(ctx).Error().Err(err).Msg("fetch products")
		}

		for _, p := range products {
			p.SellerType = c.SellerType()

			if err := s.usecase.Product.UpsertProduct(ctx, p); err != nil {
				zerolog.Ctx(ctx).Error().Err(err).Msg("upsert product")
			}
		}
	}

	return nil
}
