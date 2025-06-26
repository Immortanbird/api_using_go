package models

type RefreshTokens struct {
	UserID    string `json:"user_id" gorm:"primaryKey"`
	Token     string `json:"token" gorm:"uniqueIndex"`
	CreatedAt int64  `json:"created_at"`
}
