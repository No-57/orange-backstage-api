package store

import (
	"context"
	"orange-backstage-api/app/store/auth"

	"gorm.io/gorm"
)

type Auth struct {
	db *gorm.DB
}

func (a Auth) CreateToken(ctx context.Context, token *auth.Token) error {
	return a.db.WithContext(ctx).Create(token).Error
}
