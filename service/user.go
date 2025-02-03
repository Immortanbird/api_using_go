package service

import (
	"github.com/Immortanbird/api_using_go/model"
	"github.com/Immortanbird/api_using_go/repository"
)

func signup(user model.User) error {
	repository.InsertUser(nil, user)
	return nil
}

func login(user model.User) error {
	return nil
}
