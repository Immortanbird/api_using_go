package crud

import "github.com/Immortanbird/api_using_go/internal/models"

func CreateImage(image models.Image) error {
	result := db.Create(&image)

	return result.Error
}

func ReadImage(filter models.Image) (*models.Image, error) {
	var image models.Image
	result := db.Where(filter).First(&image)

	if result.Error != nil {
		return nil, result.Error
	}

	return &image, nil
}

func UpdateImage(filter models.Image, updateData models.Image) error {
	result := db.Model(&models.Image{}).Where(filter).Updates(updateData)

	return result.Error
}

func DeleteImage(filter models.Image) error {
	result := db.Where(filter).Delete(&models.Image{})

	return result.Error
}
