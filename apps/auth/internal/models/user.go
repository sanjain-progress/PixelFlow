package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents a registered user in the system.
type User struct {
	ID        uint           `gorm:"primaryKey"`
	Email     string         `gorm:"uniqueIndex;not null"`
	Password  string         `gorm:"not null"` // Hashed password
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
