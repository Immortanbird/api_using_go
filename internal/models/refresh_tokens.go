package models

import (
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	ID        uuid.UUID `gorm:"type:binary(16);primary_key;default:(UUID_TO_BIN(UUID(), 1))"`
	UserID    uuid.UUID `gorm:"type:binary(16);not null;constraint:OnDelete:CASCADE"`
	TokenJTI  string    `gorm:"not null;unique"`
	IsRevoked bool      `gorm:"not null;default:false"`
	IPAddress string
	UserAgent string
	ExpiresAt time.Time `gorm:"not null"`
	CreatedAt time.Time `gorm:"not null;default:now()"`
}
