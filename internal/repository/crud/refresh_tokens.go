package crud

import "github.com/Immortanbird/api_using_go/internal/models"

func CreateRefreshToken(*models.RefreshToken) error {

}

func FindRefreshToken(userID string) (string, error) {
}

func UpdateRefreshToken(userID string, newToken string) error {
}

func DeleteRefreshToken(userID string) error {

}
