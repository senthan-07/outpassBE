package models

import (
	"time"
)

type Outpass struct {
	ID           uint64    `gorm:"primaryKey"`
	StudentID    uint64    `gorm:"not null"`
	Student      Student   `gorm:"foreignKey:StudentID"`
	OutpassType  string    `gorm:"type:text"`
	Status       string    `gorm:"type:text"`
	ValidFrom    time.Time `gorm:"type:timestamptz"`
	ValidUntil   time.Time `gorm:"type:timestamptz"`
	ApprovedByID uint64    `gorm:"not null"`
	ApprovedBy   User      `gorm:"foreignKey:ApprovedByID"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type User struct {
	ID    uint64 `gorm:"primaryKey"`
	Name  string `gorm:"type:text"`
	Email string `gorm:"type:text;unique"`
}
