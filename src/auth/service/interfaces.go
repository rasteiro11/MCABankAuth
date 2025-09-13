package service

import (
	"context"

	"github.com/rasteiro11/MCABankAuth/src/auth/service/models"
	"github.com/rasteiro11/MCABankAuth/src/user/domain"
)

type (
	AuthService interface {
		Login(ctx context.Context, req *domain.User) (*models.LoginResponseDTO, error)
		Register(ctx context.Context, req *models.RegisterUserDTO) (*models.RegisterUserResponseDTO, error)
	}
)
