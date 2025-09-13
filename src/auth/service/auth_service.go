package service

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/rasteiro11/MCABankAuth/pkg/security"
	"github.com/rasteiro11/MCABankAuth/pkg/validator"
	"github.com/rasteiro11/MCABankAuth/src/auth/service/models"
	"github.com/rasteiro11/MCABankAuth/src/user/domain"
	"github.com/rasteiro11/MCABankAuth/src/user/repository"
	userService "github.com/rasteiro11/MCABankAuth/src/user/service"
)

var (
	ErrEmailTaken         = errors.New("email already taken")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type authService struct {
	userService    userService.UserService
	hasher         security.PasswordHasher
	emailValidator validator.EmailValidator
}

var _ AuthService = (*authService)(nil)

func NewAuthService(userService userService.UserService, hasher security.PasswordHasher, emailValidator validator.EmailValidator) AuthService {
	return &authService{
		userService:    userService,
		hasher:         hasher,
		emailValidator: emailValidator,
	}
}

type claims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func (s *authService) hashPassword(password string) (string, error) {
	return s.hasher.Hash(password)
}

func (s *authService) Register(ctx context.Context, req *models.RegisterUserDTO) (*models.RegisterUserResponseDTO, error) {
	if _, err := s.userService.FindOne(ctx, &domain.User{Email: req.Email}); err == nil {
		return nil, ErrEmailTaken
	} else if !errors.Is(err, repository.ErrRecordNotFound) {
		return nil, err
	}

	isEmailValid := s.emailValidator.IsValid(req.Email)
	if !isEmailValid {
		return nil, validator.ErrInvalidEmail
	}

	hash, err := s.hashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	userEntity := &domain.User{
		Email:    req.Email,
		Password: hash,
	}

	if _, err := s.userService.CreateUser(ctx, userEntity); err != nil {
		return nil, err
	}

	loginResp, err := s.Login(ctx, &domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return nil, err
	}

	return &models.RegisterUserResponseDTO{
		Token:     loginResp.Token,
		ExpiresAt: loginResp.ExpiresAt,
	}, nil
}

func (s *authService) Login(ctx context.Context, req *domain.User) (*models.LoginResponseDTO, error) {
	expiresAt := time.Now().Add(15 * time.Minute)

	userEntity, err := s.userService.FindOne(ctx, &domain.User{
		Email: req.Email,
	})
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if ok := s.hasher.Verify(req.Password, userEntity.Password); !ok {
		return nil, ErrInvalidCredentials
	}

	claims := &claims{
		UserID: userEntity.ID,
		Email:  userEntity.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return nil, err
	}

	return &models.LoginResponseDTO{
		Token:     tokenString,
		ExpiresAt: expiresAt,
	}, nil
}

func (s *authService) VerifyToken(tokenStr string) (*claims, error) {
	claims := &claims{}
	_, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	return claims, err
}
