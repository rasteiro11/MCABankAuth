package models

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type RegisterUserDTO struct {
	Email    string
	Password string
}

type LoginUserDTO struct {
	Email    string
	Password string
}

type RegisterUserResponseDTO struct {
	Token     string
	ExpiresAt time.Time
}

type LoginResponseDTO struct {
	Token     string
	ExpiresAt time.Time
}

type UserDTO struct {
	ID    uint
	Email string
}

type Claims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}
