package crud

import "github.com/Immortanbird/api_using_go/internal/models"

func CreateRefreshToken(token *models.RefreshToken) error {
	result := db.Create(token)

	return result.Error
}

func FindRefreshToken(filter models.RefreshToken) (*models.RefreshToken, error) {
	var token *models.RefreshToken
	result := db.Where(filter).First(&token)

	if result.Error != nil {
		return nil, result.Error
	}

	return token, nil
}

func RevokeRefreshToken(filter models.RefreshToken) (int64, error) {
	result := db.Where(filter).Update("is_revoked", true)

	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, nil
}

func DeleteRefreshToken(filter *models.RefreshToken) error {
	result := db.Where(filter).Delete(&models.RefreshToken{})

	return result.Error
}
