package http

import "time"

type loginRequest struct {
	Password string `json:"password" validate:"required" example:"secret123"`
	Email    string `json:"email" validate:"required" example:"rafael@example.com"`
}

type loginResponse struct {
	Token     string    `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	ExpiresAt time.Time `json:"expires_at" example:"2025-09-14T21:00:00Z"`
}

type registerRequest struct {
	Password string `json:"password" validate:"required" example:"secret123"`
	Email    string `json:"email" validate:"required" example:"rafael@example.com"`
}
