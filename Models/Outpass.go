package models

import (
	"time"

	"gorm.io/gorm"
)

type Outpass struct {
	ID           uint `gorm:"primaryKey"`
	StudentID    uint
	Student      Student
	OutpassType  string
	Status       string
	ValidFrom    time.Time
	ValidUntil   time.Time
	ApprovedByID uint
	ApprovedBy   User
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type User struct {
	gorm.Model
	Name  string
	Email string
}
