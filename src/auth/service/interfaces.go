package service

import (
	"context"

	"github.com/rasteiro11/MCABankAuth/src/user/domain"
)

type AuthService interface {
	Login(ctx context.Context, user *domain.User) (*domain.AuthSession, error)
	Register(ctx context.Context, user *domain.User) (*domain.AuthSession, error)
	VerifyToken(ctx context.Context, tokenStr string) (*domain.Claims, error)
}
