package mappers

import (
	"github.com/rasteiro11/MCABankAuth/src/user/domain"
	"github.com/rasteiro11/MCABankAuth/src/user/repository/models"
	"gorm.io/gorm"
)

func FromDomain(u *domain.User) *models.User {
	return &models.User{
		Model:    gorm.Model{ID: u.ID},
		Email:    u.Email,
		Password: u.Password,
	}
}

func ToDomain(m *models.User) *domain.User {
	return &domain.User{
		ID:       m.ID,
		Email:    m.Email,
		Password: m.Password,
	}
}
