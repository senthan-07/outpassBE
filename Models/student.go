package models

import "gorm.io/gorm"

type Student struct {
	gorm.Model
	ID        uint      `gorm:"primaryKey"`
	Name      string    `gorm:"size:255"`
	Email     string    `gorm:"unique;size:255"`
	Password  string    `gorm:"size:255"`
	Outpasses []Outpass `gorm:"foreignKey:StudentID"`
}
