package store

import (
	"fmt"
	"orange-backstage-api/app/store/account"
	"orange-backstage-api/app/store/auth"
	"orange-backstage-api/infra/util/convert"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Store struct {
	db *gorm.DB

	Account *Account
	Auth    *Auth
}

type Engine string

const (
	EnginePostgres Engine = "postgres"
	EngineMemory   Engine = "memory"
)

type Param struct {
	Engine   Engine
	Postgres PostgresCfg
	Verbose  bool
}

func New(param Param) (*Store, error) {
	var (
		db  *gorm.DB
		err error
	)

	switch param.Engine {
	case EnginePostgres:
		db, err = newPostgres(param.Verbose, param.Postgres)
	default:
		db, err = newMemory(param.Verbose)
	}
	if err != nil {
		return nil, fmt.Errorf("new db: %w", err)
	}

	if err := migrate(db); err != nil {
		return nil, fmt.Errorf("migrate: %w", err)
	}

	if err := seed(db); err != nil {
		return nil, fmt.Errorf("seed: %w", err)
	}

	return &Store{
		db: db,

		Account: &Account{db: db},
		Auth:    &Auth{db: db},
	}, nil
}

func migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(
		&account.Account{},
		&auth.Token{},
	); err != nil {
		return err
	}

	return nil
}

func seed(db *gorm.DB) error {
	rootPass := "admin"
	hashedPass, err := bcrypt.GenerateFromPassword(convert.StrToBytes(rootPass), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("generate hashed password: %w", err)
	}

	rootAcc := &account.Account{
		Email:          "admin@orange.com.tw",
		Name:           "admin",
		HashedPassword: hashedPass,
	}

	if err := db.Clauses(clause.OnConflict{
		DoNothing: true,
	}).Create(rootAcc).Error; err != nil {
		return fmt.Errorf("create root account: %w", err)
	}

	return nil
}
