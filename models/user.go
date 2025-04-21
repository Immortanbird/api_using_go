package models

import (
	"database/sql"
	"time"
)

type User struct {
	ID        uint           `json:"uid" gorm:"primarykey"`
	Name      string         `json:"name"`
	Email     sql.NullString `json:"email"`
	Age       uint8          `json:"gender"`
	Birthday  time.Time      `json:"birth"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt time.Time      `json:"deletedAt" gorm:"index"`
	// ignored      string      // fields that aren't exported are ignored
}

func InsertUser(user *User) (int64, error) {
	// Insert a new user
	result := db.Create(user)

	return result.RowsAffected, result.Error
}

func FindUsers(user *User) ([]User, error) {
	var users []User
	result := db.Find(users)

	return users, result.Error
}
