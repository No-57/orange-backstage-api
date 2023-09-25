package store

import (
	"fmt"
	"log"
	"orange-backstage-api/app/store/account"
	"orange-backstage-api/app/store/auth"
	"orange-backstage-api/infra/util/convert"
	"os"
	"time"

	"github.com/glebarez/sqlite"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	gormLogger "gorm.io/gorm/logger"
)

type Store struct {
	db *gorm.DB

	Account *Account
	Auth    *Auth
}

func New() (*Store, error) {
	db, err := gormDB()
	if err != nil {
		return nil, fmt.Errorf("gorm db: %w", err)
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

func gormDB() (*gorm.DB, error) {
	db, err := gorm.Open(
		sqlite.Open("file::memory:?cache=shared"),
		&gorm.Config{
			PrepareStmt: true,
			Logger: gormLogger.New(
				log.New(os.Stdout, "\r\n", log.LstdFlags),
				gormLogger.Config{
					SlowThreshold:             time.Second,
					LogLevel:                  gormLogger.Info,
					IgnoreRecordNotFoundError: true,
					Colorful:                  true,
				},
			),
			CreateBatchSize:                          1000,
			DisableForeignKeyConstraintWhenMigrating: true,
			SkipDefaultTransaction:                   true,
			QueryFields:                              true,
		},
	)
	if err != nil {
		return nil, err
	}

	return db, nil
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
