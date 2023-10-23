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

func (u *User) FindByEmail(db *gorm.DB, username string) (*User, error) {
	var (
		err  error
		user User
	)

	err = db.Debug().Model(User{}).Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil

}

func (u *User) FindByID(db *gorm.DB, userID string) (*User, error) {
	var (
		err  error
		user User
	)

	err = db.Debug().Model(User{}).Where("id = ?", userID).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil

}
