package repository

import (
	"github.com/Immortanbird/api_using_go/model"
	"gorm.io/gorm"
)

func InsertUser(db *gorm.DB, user model.User) error {
	// Insert a new user
	db.Create(&user)
	return nil
}
