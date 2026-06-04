package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Post represents the structure for a blog post in the database.
type Post struct {
	ID     uuid.UUID `json:"id" gorm:"type:binary(16);primaryKey;default:(UUID_TO_BIN(UUID(), 1))"`
	UserID uuid.UUID `json:"user_id" gorm:"type:binary(16);not null"`

	Title string `json:"title" gorm:"type:varchar(255);not null" validate:"required,min=1,max=255"`

	// A URL-friendly version of the title.
	Slug string `json:"slug" gorm:"type:varchar(255);not null;unique" validate:"required,min=1,max=255"`

	Excerpt string `json:"excerpt" gorm:"type:text"`

	ContentMarkdown string `json:"content_markdown" gorm:"type:longtext;not null" validate:"required"`
	ContentHTML     string `json:"content_html" gorm:"type:longtext;not null" validate:"required"`

	// Cover image
	CoverImageURL string `json:"cover_image_url" gorm:"type:varchar(500)"`

	// The publication status, e.g., "draft" or "published".
	Status string `json:"status" gorm:"type:enum('draft','published','archived','private');not null;default:'draft'" validate:"oneof=draft published archived private"`

	// Publishing control
	PublishedAt *time.Time `json:"published_at" gorm:"type:timestamp null"`
	ScheduledAt *time.Time `json:"scheduled_at" gorm:"type:timestamp null"`

	// Engagement metrics
	ViewCount    int64 `json:"view_count" gorm:"type:bigint;default:0"`
	LikeCount    int64 `json:"like_count" gorm:"type:bigint;default:0"`
	CommentCount int64 `json:"comment_count" gorm:"type:bigint;default:0"`

	// Reading time estimation (in minutes)
	ReadingTime int `json:"reading_time" gorm:"type:int;default:0"`

	// Soft delete support
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"type:timestamp null;index"`

	CreatedAt time.Time `json:"created_at" gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `json:"updated_at" gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`

	// Define the relationship to the Users model.
	Author Users `json:"author" gorm:"foreignKey:UserID;references:ID"`

	// Many-to-many relationship with tags
	Tags []Tag `json:"tags" gorm:"many2many:posts_tags;"`

	// One-to-many relationship with comments
	Comments []Comment `json:"comments" gorm:"foreignKey:PostID;references:ID"`
}

// Tag represents a tag that can be associated with posts
type Tag struct {
	ID          uuid.UUID `json:"id" gorm:"type:binary(16);primaryKey;default:(UUID_TO_BIN(UUID(), 1))"`
	Name        string    `json:"name" gorm:"type:varchar(100);not null;unique" validate:"required,min=1,max=100"`
	Slug        string    `json:"slug" gorm:"type:varchar(100);not null;unique" validate:"required,min=1,max=100"`
	Description string    `json:"description" gorm:"type:text"`
	CreatedAt   time.Time `json:"created_at" gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`

	// Many-to-many relationship with posts
	Posts []Post `json:"posts" gorm:"many2many:posts_tags;"`
}

// Comment represents a comment on a post
type Comment struct {
	ID         uuid.UUID  `json:"id" gorm:"type:binary(16);primaryKey;default:(UUID_TO_BIN(UUID(), 1))"`
	PostID     uuid.UUID  `json:"post_id" gorm:"type:binary(16);not null"`
	UserID     *uuid.UUID `json:"user_id" gorm:"type:binary(16);null"`   // Allow anonymous comments
	ParentID   *uuid.UUID `json:"parent_id" gorm:"type:binary(16);null"` // For nested comments
	AuthorName string     `json:"author_name" gorm:"type:varchar(100);not null" validate:"required,min=1,max=100"`
	Content    string     `json:"content" gorm:"type:text;not null" validate:"required"`
	CreatedAt  time.Time  `json:"created_at" gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt  time.Time  `json:"updated_at" gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`

	// Relationships
	Post    Post      `json:"post" gorm:"foreignKey:PostID;references:ID"`
	User    *Users    `json:"user" gorm:"foreignKey:UserID;references:ID"`
	Parent  *Comment  `json:"parent" gorm:"foreignKey:ParentID;references:ID"`
	Replies []Comment `json:"replies" gorm:"foreignKey:ParentID;references:ID"`
}

// PostTag represents the junction table for posts and tags
type PostTag struct {
	ID        uuid.UUID `json:"id" gorm:"type:binary(16);primaryKey;default:(UUID_TO_BIN(UUID(), 1))"`
	PostID    uuid.UUID `json:"post_id" gorm:"type:binary(16);not null"`
	TagID     uuid.UUID `json:"tag_id" gorm:"type:binary(16);not null"`
	CreatedAt time.Time `json:"created_at" gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
}
