package models

import (
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID    uuid.UUID `gorm:"type:uuid;not null"`
	TokenJTI  string    `gorm:"not null;unique"`
	IsRevoked bool      `gorm:"not null;default:false"`
	IPAddress string
	UserAgent string
	ExpiresAt time.Time `gorm:"not null"`
	CreatedAt time.Time `gorm:"not null;default:now()"`
}
