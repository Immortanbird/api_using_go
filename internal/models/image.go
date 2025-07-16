package models

import (
	"time"

	"github.com/google/uuid"
)

type Image struct {
	ID               uuid.UUID `gorm:"type:binary(16);primaryKey;default:(UUID_TO_BIN(UUID(), 1))"`
	UserID           uuid.UUID `gorm:"type:binary(16);not null;constraint:OnDelete:CASCADE"`
	OriginalFilename string    `gorm:"type:text;not null"`
	HashedFilename   string    `gorm:"type:text;not null;unique"`
	ImageURL         string    `gorm:"type:text;not null;unique"`
	DownloadURL      string    `gorm:"type:text;not null;unique"`
	SizeBytes        int64     `gorm:"type:bigint;not null"`
	MimeType         string    `gorm:"type:varchar(128);not null"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
