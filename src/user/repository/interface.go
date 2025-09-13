package repository

import (
	"context"

	"github.com/rasteiro11/MCABankAuth/src/user/domain"
)

type (
	Repository interface {
		FindOne(ctx context.Context, user *domain.User) (*domain.User, error)
		Create(ctx context.Context, user *domain.User) (*domain.User, error)
		FindOneByEmail(ctx context.Context, email string) (*domain.User, error)
	}
)
