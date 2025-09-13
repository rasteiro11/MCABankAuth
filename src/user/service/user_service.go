package service

import (
	"context"

	"github.com/rasteiro11/MCABankAuth/src/user/domain"
	"github.com/rasteiro11/MCABankAuth/src/user/repository"
)

type userService struct {
	repo repository.Repository
}

var _ UserService = (*userService)(nil)

func NewUserService(repo repository.Repository) UserService {
	return &userService{repo: repo}
}

func (s *userService) GetUserByID(ctx context.Context, id uint) (*domain.User, error) {
	return s.repo.FindOne(ctx, &domain.User{ID: id})
}

func (s *userService) GetUserByDocument(ctx context.Context, doc string) (*domain.User, error) {
	return s.repo.FindOne(ctx, &domain.User{Document: doc})
}

func (s *userService) FindOne(ctx context.Context, user *domain.User) (*domain.User, error) {
	return s.repo.FindOne(ctx, user)
}

func (s *userService) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	return s.repo.Create(ctx, user)
}
