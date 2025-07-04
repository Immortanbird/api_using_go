package models

import (
	"time"

	"github.com/google/uuid"
)

type Users struct {
	ID        uuid.UUID `json:"uid" gorm:"binary(16);primaryKey;default:(UUID_TO_BIN(UUID(), 1))"`
	Email     string    `json:"email" gorm:"not null;unique"`
	Username  string    `json:"username" gorm:"varchar(50);not null"`
	Password  string    `json:"-" gorm:"varchar(255);not null"`
	Age       uint8     `json:"age"`
	Birthday  time.Time `json:"birth"`
	Verified  bool      `gorm:"verified"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
