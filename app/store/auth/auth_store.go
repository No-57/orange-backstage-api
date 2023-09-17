package auth

import "gorm.io/gorm"

type Store struct {
	db *gorm.DB
}
