package models

import (
	"time"

	"gorm.io/gorm"
)

type Address struct {
	ID        string `gorm:"size:36; not null; uniqueIndex;primary_key"`
	User      User
	UserID    string `gorm:"size:36;index"`
	Name      string `gorm:"size:100"`
	IsPrimary bool
	Address1  string `gorm:"size:255"`
	Address2  string `gorm:"size:255"`
	Phone     string `gorm:"size:50"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}
