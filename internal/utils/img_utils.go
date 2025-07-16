package utils

import (
	"crypto/sha256"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"time"
)

func HashFileName(fileHeader *multipart.FileHeader) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", err
	}

	hashString := fmt.Sprintf("%x", hasher.Sum(nil))
	extension := filepath.Ext(fileHeader.Filename)
	timestamp := time.Now().Format("20010101121212")
	return fmt.Sprintf("%s-%s%s", timestamp, hashString[:16], extension), nil
}

func IsValidImageType(fileHeader *multipart.FileHeader) (bool, string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return false, "", err
	}
	defer file.Close()

	// Read the first 512 bytes to detect the content type.
	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil && err != io.EOF {
		return false, "", err
	}

	contentType := http.DetectContentType(buffer)
	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/jpg":  true,
		"image/png":  true,
		// "image/gif":  true,
		// "image/webp": true,
	}

	return allowedTypes[contentType], contentType, nil
}
