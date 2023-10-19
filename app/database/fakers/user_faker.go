package fakers

import (
	"simpleWebCart/app/models"
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func UserFaker(db *gorm.DB) *models.User {
	return &models.User{
		ID:        uuid.New().String(),
		Name:      faker.Name(),
		Username:  faker.Username(),
		Password:  faker.Password(),
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
		DeletedAt: gorm.DeletedAt{},
	}
}
