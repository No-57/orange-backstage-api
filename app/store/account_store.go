package store

import (
	"context"
	"orange-backstage-api/app/store/account"

	"gorm.io/gorm"
)

type Account struct {
	db *gorm.DB
}

func (a *Account) GetByEmailOrName(ctx context.Context, val string) (*account.Account, error) {
	var record *account.Account
	if err := a.db.WithContext(ctx).
		Where(account.Account{Email: val}).
		Or(account.Account{Name: val}).
		Take(&record).Error; err != nil {
		return nil, err
	}
	return record, nil
}
