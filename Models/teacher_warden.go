package models

import "gorm.io/gorm"

type Warden struct {
	gorm.Model
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"size:255"`
	Email     string    `gorm:"unique;size:255"`
	Password  string    `gorm:"size:255"`
	Outpasses []Outpass `gorm:"foreignKey:ApprovedByID"` // Outpasses approved by this warden
}

type Teacher struct {
	gorm.Model
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:255"`
	Email     string    `gorm:"unique;size:255"`
	Password  string    `gorm:"size:255"`
	Outpasses []Outpass `gorm:"foreignKey:ApprovedByID"` // Outpasses approved by this teacher
}
