package models

type LoginRequest struct {
	Email    string
	Username string
	Password string
}

type RegisterRequest struct {
	Email    string
	Username string
	Password string
}
