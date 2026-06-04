package handler

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Immortanbird/api_using_go/internal/models"
	"github.com/Immortanbird/api_using_go/internal/repository/crud"
	"github.com/Immortanbird/api_using_go/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// POST	/image/upload
func (h *Handler) UploadImage(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	fileHeader, err := c.FormFile("image")
	if err != nil {
		zap.L().Warn(
			"No images found",
			zap.Error(err),
			zap.String("client_ip", c.ClientIP()),
			zap.String("path", c.Request.URL.Path),
			zap.String("userID", userID.String()),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Image file is required in 'image' field."})
		return
	}

	valid, contentType, err := utils.IsValidImageType(fileHeader)
	if !valid || err != nil {
		if err != nil {
			zap.L().Error(
				"Failed to validate image type",
				zap.Error(err),
				zap.String("client_ip", c.ClientIP()),
				zap.String("path", c.Request.URL.Path),
				zap.String("userID", userID.String()),
			)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal serve error."})
		} else {
			c.JSON(http.StatusUnsupportedMediaType, gin.H{"error": "Invalid image type."})
		}
		return
	}

	hashedImageName, err := utils.HashFileName(fileHeader)
	if err != nil {
		zap.L().Error(
			"Failed to generate hashed image name.",
			zap.Error(err),
			zap.String("client_ip", c.ClientIP()),
			zap.String("path", c.Request.URL.Path),
			zap.String("userID", userID.String()),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error."})
		return
	}

	file, _ := fileHeader.Open()
	defer file.Close()
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		zap.L().Error(
			"Failed to read file content.",
			zap.Error(err),
			zap.String("client_ip", c.ClientIP()),
			zap.String("path", c.Request.URL.Path),
			zap.String("userID", userID.String()),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not read file content"})
		return
	}

	commitMessage := fmt.Sprintf("feat: Add image %s", hashedImageName)

	apiUrl := fmt.Sprintf(h.Config.PicBed.GitHubUrlFormat, hashedImageName)
	encodedContent := base64.StdEncoding.EncodeToString(fileBytes)

	commitData, _ := json.Marshal(map[string]string{
		"message": commitMessage,
		"content": encodedContent,
	})

	client := &http.Client{Timeout: 10 * time.Second}

	request, _ := http.NewRequest(http.MethodPut, apiUrl, bytes.NewBuffer(commitData))
	request.Header.Set("Authorization", "Bearer "+h.Config.PicBed.GitHubToken)
	request.Header.Set("Accept", "application/vnd.github+json")
	response, err := client.Do(request)
	if err != nil || response.StatusCode >= 300 {
		zap.L().Error(
			"Failed to upload image to GitHub.",
			zap.Error(err),
			zap.String("client_ip", c.ClientIP()),
			zap.String("path", c.Request.URL.Path),
			zap.String("userID", userID.String()),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error."})
		return
	}
	defer response.Body.Close()

	// Get image info
	request, _ = http.NewRequest(http.MethodGet, apiUrl, nil)
	response, err = client.Do(request)
	if err != nil || response.StatusCode != http.StatusOK {
		zap.L().Error(
			"Failed to get image info.",
			zap.Error(err),
			zap.String("client_ip", c.ClientIP()),
			zap.String("path", c.Request.URL.Path),
			zap.String("userID", userID.String()),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error."})
		return
	}
	defer response.Body.Close()

	bodyBytes, _ := io.ReadAll(response.Body)
	jsonBody := make(map[string]interface{})
	if err := json.Unmarshal(bodyBytes, &jsonBody); err != nil {
		zap.L().Error(
			"Failed to unmarshal image info.",
			zap.Error(err),
			zap.String("client_ip", c.ClientIP()),
			zap.String("path", c.Request.URL.Path),
			zap.String("userID", userID.String()),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error."})
		return
	}

	image := models.Image{
		UserID:           userID,
		OriginalFilename: fileHeader.Filename,
		HashedFilename:   hashedImageName,
		ImageURL:         jsonBody["url"].(string),
		DownloadURL:      jsonBody["download_url"].(string),
		SizeBytes:        jsonBody["size"].(int64),
		MimeType:         contentType,
	}

	if err := crud.CreateImage(image); err != nil {
		zap.L().Error(
			"Failed to create image.",
			zap.Error(err),
			zap.String("client_ip", c.ClientIP()),
			zap.String("path", c.Request.URL.Path),
			zap.String("userID", userID.String()),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error."})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"image_name": hashedImageName})
}

// GET /image/download
func (h *Handler) DownloadImage(c *gin.Context) {

}

// DELETE /image/delete
func (h *Handler) DeleteImage(c *gin.Context) {

}
