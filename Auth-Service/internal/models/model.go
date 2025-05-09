package models

import (
	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primaryKey"`
	Email     string         `gorm:"unique;not null"`
	Password  string         `gorm:"not null"`
	isBlocked bool           `gorm:"is_blocked"`
	Role      string         `gorm:"not null;default:user"` // "user" or "admin"
	CreatedAt int64          `gorm:"autoCreateTime"`
	UpdatedAt int64          `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
