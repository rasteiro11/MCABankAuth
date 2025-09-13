package http

import "time"

type loginRequest struct {
	Password string `json:"password" validate:"required"`
	Email    string `json:"email" validate:"required"`
}

type loginResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

type registerRequest struct {
	Password string `json:"password" validate:"required"`
	Email    string `json:"email" validate:"required"`
}
