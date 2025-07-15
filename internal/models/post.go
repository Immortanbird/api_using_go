package models

import (
	"time"

	"github.com/google/uuid"
)

// Post represents the structure for a blog post in the database.
type Post struct {
	ID uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`

	UserID uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`

	Title string `json:"title" gorm:"type:varchar(255);not null"`

	// A URL-friendly version of the title.
	Slug string `json:"slug" gorm:"type:text;not null;unique"`

	Excerpt string `json:"excerpt" gorm:"type:text"`

	Content string `json:"content" gorm:"type:text;not null"`

	CoverImageURL string `json:"cover_image_url" gorm:"type:varchar(255)"`

	// The publication status, e.g., "draft" or "published".
	Status string `json:"status" gorm:"type:varchar(50);not null;default:'draft'"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Define the relationship to the Users model.
	// This tells GORM that a Post belongs to a User.
	// The `foreignKey` and `references` tags ensure the relationship is correctly mapped.
	Author Users `json:"author" gorm:"foreignKey:UserID;references:ID"`
}
