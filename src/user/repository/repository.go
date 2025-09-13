package repository

import (
	"context"

	"github.com/rasteiro11/MCABankAuth/src/user/domain"
	"github.com/rasteiro11/MCABankAuth/src/user/repository/models"
	"github.com/rasteiro11/MCABankAuth/src/user/repository/models/mappers"
	"github.com/rasteiro11/PogCore/pkg/database"
	"github.com/rasteiro11/PogCore/pkg/logger"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

var _ Repository = (*repository)(nil)

var ErrRecordNotFound = gorm.ErrRecordNotFound

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindOne(ctx context.Context, u *domain.User) (*domain.User, error) {
	res := &models.User{}

	if err := r.db.Where(mappers.FromDomain(u)).Take(res).Error; err != nil {
		logger.Of(ctx).Errorf("[user.repository.FindOne] db.Take() returned error: %+v\n", err)
		return nil, err
	}

	return mappers.ToDomain(res), nil
}

func (r *repository) Create(ctx context.Context, u *domain.User) (*domain.User, error) {
	model := mappers.FromDomain(u)

	tx, err := database.FromContext(ctx)
	if err != nil {
		tx = r.db
	}

	if err := tx.Debug().Create(model).Error; err != nil {
		logger.Of(ctx).Errorf("[user.repository.Create] db.Create() returned error: %+v\n", err)
		return nil, err
	}

	u.ID = model.ID
	return u, nil
}

func (r *repository) FindOneByEmail(ctx context.Context, email string) (*domain.User, error) {
	res := &models.User{}
	if err := r.db.Where(&models.User{Email: email}).Take(res).Error; err != nil {
		logger.Of(ctx).Errorf("[user.repository.FindOneByEmail] db.Take() returned error: %+v\n", err)
		return nil, err
	}

	return mappers.ToDomain(res), nil
}
