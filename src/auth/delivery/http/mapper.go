package http

import (
	"github.com/rasteiro11/MCABankAuth/src/user/domain"
)

func MapLoginRequestToUser(req *loginRequest) *domain.User {
	return &domain.User{
		Email:    req.Email,
		Password: req.Password,
	}
}

func MapUserLoginResponseToHTTP(resp *domain.AuthSession) *authSession {
	return &authSession{
		Token:     resp.Token,
		ExpiresAt: resp.ExpiresAt,
	}
}

func MapRegisterRequestToDTO(req *registerRequest) *domain.User {
	return &domain.User{
		Email:    req.Email,
		Password: req.Password,
	}
}

func MapUserRegisterResponseToHTTP(resp *domain.AuthSession) *authSession {
	return &authSession{
		Token:     resp.Token,
		ExpiresAt: resp.ExpiresAt,
	}
}
