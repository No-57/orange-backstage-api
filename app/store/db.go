package store

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	gormLogger "gorm.io/gorm/logger"
)

type PostgresCfg struct {
	Host     string
	User     string
	Password string
	DBName   string
	Port     string
	SSLMode  string
	TimeZone string
}

func newPostgres(verbose bool, p PostgresCfg) (*gorm.DB, error) {
	logLevel := gormLogger.Error
	if verbose {
		logLevel = gormLogger.Info
	}

	dsn := postgresDSN(p)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		PrepareStmt: true,
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold:             time.Second,
				LogLevel:                  logLevel,
				IgnoreRecordNotFoundError: true,
				Colorful:                  true,
			},
		),
		CreateBatchSize:                          1000,
		SkipDefaultTransaction:                   true,
		DisableForeignKeyConstraintWhenMigrating: true,
		QueryFields:                              true,
	})
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get sql db: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}

func postgresDSN(p PostgresCfg) string {
	const dsnFormat = "host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s"

	return fmt.Sprintf(
		dsnFormat,
		p.Host,
		p.User,
		p.Password,
		p.DBName,
		p.Port,
		p.SSLMode,
		p.TimeZone,
	)
}

func newMemory(verbose bool) (*gorm.DB, error) {
	logLevel := gormLogger.Error
	if verbose {
		logLevel = gormLogger.Info
	}

	db, err := gorm.Open(
		sqlite.Open("file::memory:?cache=shared"),
		&gorm.Config{
			PrepareStmt: true,
			Logger: gormLogger.New(
				log.New(os.Stdout, "\r\n", log.LstdFlags),
				gormLogger.Config{
					SlowThreshold:             time.Second,
					LogLevel:                  logLevel,
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
