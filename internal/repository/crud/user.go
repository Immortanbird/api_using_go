package crud

import (
	"github.com/Immortanbird/api_using_go/internal/models"
)

func CreateUser(user *models.Users) error {
	result := db.Create(user)

	return result.Error
}

func FindTheUser(filter models.Users) (*models.Users, error) {
	var user models.Users
	result := db.Where(filter).First(&user)

	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func FindUsers(filter models.Users) ([]models.Users, error) {
	var users []models.Users
	result := db.Where(filter).First(&users)

	if result.Error != nil {
		return []models.Users{}, result.Error
	}

	return users, nil
}

func UpdateUser() {

}

func DeleteUser() {
}
