package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        string `gorm:"size:36; not null; uniqueIndex;primary_key"`
	Adresses  []Address
	Name      string `gorm:"size:100; not null"`
	Username  string `gorm:"size:100; not null; uniqueIndex"`
	Password  string `gorm:"size:100; not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}
