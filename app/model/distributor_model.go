package model

import "time"

type Distributor struct {
	ID          uint64    `gorm:"column:id;type:serial;primaryKey"`
	Name        string    `gorm:"column:name;NOT NULL"`
	Description string    `gorm:"column:description;NOT NULL"`
	BrandImage  string    `gorm:"column:brand_image;NOT NULL"`
	CreatedAt   time.Time `gorm:"column:created_date;NOT NULL"`
	UpdateAt    time.Time `gorm:"column:update_date;NOT NULL"`
}
