package service

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"github.com/rasteiro11/MCABankAuth/pkg/security"
	"github.com/rasteiro11/MCABankAuth/pkg/token"
	"github.com/rasteiro11/MCABankAuth/pkg/validator"
	"github.com/rasteiro11/MCABankAuth/src/user/domain"
	"github.com/rasteiro11/MCABankAuth/src/user/repository"
	userService "github.com/rasteiro11/MCABankAuth/src/user/service"
	"time"
)

var (
	ErrEmailTaken         = errors.New("email already taken")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrSignatureInvalid   = errors.New("error invalid signature")
)

type authService struct {
	userService    userService.UserService
	hasher         security.PasswordHasher
	emailValidator validator.EmailValidator
	jwtService     token.JWTService[domain.Claims]
}

var _ AuthService = (*authService)(nil)

func NewAuthService(userService userService.UserService, hasher security.PasswordHasher, emailValidator validator.EmailValidator, jwtService token.JWTService[domain.Claims]) AuthService {
	return &authService{
		userService:    userService,
		hasher:         hasher,
		emailValidator: emailValidator,
		jwtService:     jwtService,
	}
}

func (s *authService) hashPassword(password string) (string, error) {
	return s.hasher.Hash(password)
}

func (s *authService) Register(ctx context.Context, req *domain.User) (*domain.AuthSession, error) {
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

	return &domain.AuthSession{
		Token:     loginResp.Token,
		ExpiresAt: loginResp.ExpiresAt,
	}, nil
}

func (s *authService) Login(ctx context.Context, req *domain.User) (*domain.AuthSession, error) {
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

	claims := domain.Claims{
		UserID: userEntity.ID,
		Email:  userEntity.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}

	token, err := s.jwtService.CreateToken(claims)
	if err != nil {
		return nil, err
	}

	return &domain.AuthSession{
		Token:     token,
		ExpiresAt: expiresAt,
	}, nil
}

func (s *authService) VerifyToken(ctx context.Context, tokenStr string) (*domain.Claims, error) {
	claims, err := s.jwtService.ParseToken(tokenStr)
	if err != nil {
		return nil, err
	}

	return &claims, nil
}
