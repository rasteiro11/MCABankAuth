package service

import (
	"context"

	"github.com/rasteiro11/MCABankAuth/src/user/domain"
)

type UserService interface {
	GetUserByID(ctx context.Context, id uint) (*domain.User, error)
	GetUserByDocument(ctx context.Context, document string) (*domain.User, error)
	FindOne(ctx context.Context, user *domain.User) (*domain.User, error)
	CreateUser(ctx context.Context, user *domain.User) (*domain.User, error)
}
