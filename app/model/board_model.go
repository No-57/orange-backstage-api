package model

import "time"

type Board struct {
	ID         uint64    `gorm:"column:id;type:serial;primaryKey"`
	Code       string    `gorm:"column:code;NOT NULL"`
	ImageURL   string    `gorm:"column:image_url;NOT NULL"`
	ActionType string    `gorm:"column:action_type;NOT NULL"`
	Action     string    `gorm:"column:action;NOT NULL"`
	CreatedAt  time.Time `gorm:"column:created_at;NOT NULL"`
	UpdatedAt  time.Time `gorm:"column:updated_at;NOT NULL"`
}
