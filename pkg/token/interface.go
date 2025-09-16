package token

import "github.com/golang-jwt/jwt/v4"

type JWTService[T jwt.Claims] interface {
	CreateToken(claims T) (string, error)
	ParseToken(tokenStr string) (T, error)
	IsValid(tokenStr string) bool
}
