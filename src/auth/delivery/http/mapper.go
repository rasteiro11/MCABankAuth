package http

import (
	"github.com/rasteiro11/MCABankAuth/src/auth/service/models"
	"github.com/rasteiro11/MCABankAuth/src/user/domain"
)

func MapLoginRequestToUser(req *loginRequest) *domain.User {
	return &domain.User{
		Email:    req.Email,
		Password: req.Password,
	}
}

func MapUserLoginResponseToHTTP(resp *models.LoginResponseDTO) *loginResponse {
	return &loginResponse{
		Token:     resp.Token,
		ExpiresAt: resp.ExpiresAt,
	}
}

func MapRegisterRequestToDTO(req *registerRequest) *models.RegisterUserDTO {
	return &models.RegisterUserDTO{
		Email:    req.Email,
		Password: req.Password,
	}
}

func MapUserRegisterResponseToHTTP(resp *models.RegisterUserResponseDTO) *loginResponse {
	return &loginResponse{
		Token:     resp.Token,
		ExpiresAt: resp.ExpiresAt,
	}
}
