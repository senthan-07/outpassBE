package models

import (
	"time"

	"gorm.io/gorm"
)

type Outpass struct {
	ID           uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	StudentID    uint64         `gorm:"not null" json:"student_id"`
	OutpassType  string         `gorm:"not null" json:"outpass_type"`
	Status       string         `gorm:"default:Pending" json:"status"`
	ValidFrom    time.Time      `gorm:"not null" json:"valid_from"`
	ValidUntil   time.Time      `gorm:"not null" json:"valid_until"`
	ApprovedByID *uint64        `json:"approved_by_id"`                // Can be warden or teacher
	ApproverType string         `gorm:"not null" json:"approver_type"` // "warden" or "teacher"
	CreatedAt    time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

type OutpassResponse struct {
	ID           uint      `json:"id"`
	StudentID    uint      `json:"student_id"`
	OutpassType  string    `json:"outpass_type"`
	Status       string    `json:"status"`
	ValidFrom    time.Time `json:"valid_from"`
	ValidUntil   time.Time `json:"valid_until"`
	ApprovedByID *uint     `json:"approved_by_id,omitempty"` // Omitting null values
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type User struct {
	ID    uint64 `gorm:"primaryKey"`
	Name  string `gorm:"type:text"`
	Email string `gorm:"type:text;unique"`
}

type Notification struct {
	ID      uint64 `gorm:"primaryKey"`
	UserID  uint64
	Message string
	Read    bool
}
