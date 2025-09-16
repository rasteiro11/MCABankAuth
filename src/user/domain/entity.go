package domain

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type User struct {
	ID       uint
	Email    string
	Password string
}

type AuthSession struct {
	Token     string
	ExpiresAt time.Time
}

type Claims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}
