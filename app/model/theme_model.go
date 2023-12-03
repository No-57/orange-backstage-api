package model

import "time"

type Theme struct {
	ID        uint64    `gorm:"column:id;type:serial;primaryKey"`
	Code      string    `gorm:"column:code;NOT NULL"`
	Type      string    `gorm:"column:type;NOT NULL"`
	Disable   bool      `gorm:"column:disable;NOT NULL"`
	CreatedAt time.Time `gorm:"column:created_date;NOT NULL"`
	UpdatedAt time.Time `gorm:"column:updated_date;NOT NULL"`
}
