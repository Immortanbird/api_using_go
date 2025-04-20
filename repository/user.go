package repository

import (
	"github.com/Immortanbird/api_using_go/model"
)

func InsertUser(user model.User) error {
	// Insert a new user
	db.Create(&user)
	return nil
}
