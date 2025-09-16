package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type jwtService[T jwt.Claims] struct {
	secret     []byte
	expiration time.Duration
}

func NewJWTService[T jwt.Claims](secret string, expiration time.Duration) JWTService[T] {
	return &jwtService[T]{
		secret:     []byte(secret),
		expiration: expiration,
	}
}

func (s *jwtService[T]) CreateToken(claims T) (string, error) {
	switch c := any(&claims).(type) {
	case *struct {
		jwt.RegisteredClaims
	}:
		if c.ExpiresAt == nil {
			c.ExpiresAt = jwt.NewNumericDate(time.Now().Add(s.expiration))
		}
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secret)
}

func (s *jwtService[T]) ParseToken(tokenStr string) (T, error) {
	var zero T

	claimsPtr := new(T)
	claimsAsJWT, ok := any(claimsPtr).(jwt.Claims)
	if !ok {
		return zero, errors.New("*T does not implement jwt.Claims")
	}
	token, err := jwt.ParseWithClaims(tokenStr, claimsAsJWT, func(t *jwt.Token) (interface{}, error) {
		return s.secret, nil
	})
	if err != nil {
		return zero, err
	}

	if !token.Valid {
		return zero, ErrSignatureInvalid
	}

	return *claimsPtr, nil
}

func (s *jwtService[T]) IsValid(tokenStr string) bool {
	_, err := s.ParseToken(tokenStr)
	return err == nil
}
