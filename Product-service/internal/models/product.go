package models

import (
	"gorm.io/gorm"
)

type Product struct {
	ID        uint           `gorm:"primaryKey"`
	Name      string         `gorm:"not null"`
	Price     float32        `gorm:"not null"`
	Quantity  int32          `gorm:"not null"`
	CreatedAt int64          `gorm:"autoCreateTime"`
	UpdatedAt int64          `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
