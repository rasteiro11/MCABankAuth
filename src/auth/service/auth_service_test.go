package service_test

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	securityMocks "github.com/rasteiro11/MCABankAuth/pkg/security/mocks"
	"github.com/rasteiro11/MCABankAuth/pkg/token"
	validatorImpl "github.com/rasteiro11/MCABankAuth/pkg/validator"
	validatorMocks "github.com/rasteiro11/MCABankAuth/pkg/validator/mocks"
	authService "github.com/rasteiro11/MCABankAuth/src/auth/service"
	"github.com/rasteiro11/MCABankAuth/src/user/domain"
	userMocks "github.com/rasteiro11/MCABankAuth/src/user/mocks"
	repo "github.com/rasteiro11/MCABankAuth/src/user/repository"
)

func TestAuthService(t *testing.T) {
	ctx := context.Background()
	os.Setenv("JWT_SECRET", "testsecret")

	userSvc := new(userMocks.UserService)
	hasher := new(securityMocks.PasswordHasher)
	validator := new(validatorMocks.EmailValidator)
	jwtSvc := token.NewJWTService[domain.Claims]("testsecret", 15*time.Minute)

	svc := authService.NewAuthService(userSvc, hasher, validator, jwtSvc)

	t.Run("Register - email already taken", func(t *testing.T) {
		userSvc.On("FindOne", ctx, mock.Anything).Return(&domain.User{}, nil).Once()

		_, err := svc.Register(ctx, &domain.User{Email: "taken@example.com", Password: "pwd"})
		assert.ErrorIs(t, err, authService.ErrEmailTaken)

		userSvc.AssertExpectations(t)
	})

	t.Run("Register - invalid email", func(t *testing.T) {
		userSvc.On("FindOne", ctx, mock.Anything).Return(nil, repo.ErrRecordNotFound).Once()
		validator.On("IsValid", "bad-email").Return(false).Once()

		_, err := svc.Register(ctx, &domain.User{Email: "bad-email", Password: "pwd"})
		assert.ErrorIs(t, err, validatorImpl.ErrInvalidEmail)

		userSvc.AssertExpectations(t)
		validator.AssertExpectations(t)
	})

	t.Run("Register - success", func(t *testing.T) {
		req := &domain.User{Email: "new@example.com", Password: "pwd"}

		userSvc.On("FindOne", ctx, mock.Anything).Return(nil, repo.ErrRecordNotFound).Once()
		validator.On("IsValid", req.Email).Return(true).Once()
		hasher.On("Hash", req.Password).Return("hashedpwd", nil).Once()

		userSvc.On("CreateUser", ctx, mock.Anything).Return(&domain.User{
			ID:       1,
			Email:    req.Email,
			Password: "hashedpwd",
		}, nil).Once()

		userSvc.On("FindOne", ctx, mock.Anything).Return(&domain.User{
			ID:       1,
			Email:    req.Email,
			Password: "hashedpwd",
		}, nil).Once()
		hasher.On("Verify", req.Password, "hashedpwd").Return(true).Once()

		resp, err := svc.Register(ctx, req)
		assert.NoError(t, err)
		assert.NotEmpty(t, resp.Token)
		assert.WithinDuration(t, time.Now().Add(15*time.Minute), resp.ExpiresAt, 2*time.Second)

		userSvc.AssertExpectations(t)
		validator.AssertExpectations(t)
		hasher.AssertExpectations(t)
	})

	t.Run("Login - invalid user", func(t *testing.T) {
		userSvc.On("FindOne", ctx, mock.Anything).Return(nil, errors.New("not found")).Once()

		_, err := svc.Login(ctx, &domain.User{Email: "x@example.com", Password: "pwd"})
		assert.ErrorIs(t, err, authService.ErrInvalidCredentials)

		userSvc.AssertExpectations(t)
	})

	t.Run("Login - wrong password", func(t *testing.T) {
		userSvc.On("FindOne", ctx, mock.Anything).Return(&domain.User{
			ID:       1,
			Email:    "x@example.com",
			Password: "hashed",
		}, nil).Once()

		hasher.On("Verify", "pwd", "hashed").Return(false).Once()

		_, err := svc.Login(ctx, &domain.User{Email: "x@example.com", Password: "pwd"})
		assert.ErrorIs(t, err, authService.ErrInvalidCredentials)

		userSvc.AssertExpectations(t)
		hasher.AssertExpectations(t)
	})

	t.Run("Login - success", func(t *testing.T) {
		userSvc.On("FindOne", ctx, mock.Anything).Return(&domain.User{
			ID:       1,
			Email:    "x@example.com",
			Password: "hashed",
		}, nil).Once()
		hasher.On("Verify", "pwd", "hashed").Return(true).Once()

		resp, err := svc.Login(ctx, &domain.User{Email: "x@example.com", Password: "pwd"})
		assert.NoError(t, err)
		assert.NotEmpty(t, resp.Token)
		assert.WithinDuration(t, time.Now().Add(15*time.Minute), resp.ExpiresAt, 2*time.Second)

		userSvc.AssertExpectations(t)
		hasher.AssertExpectations(t)
	})

	t.Run("VerifyToken - success", func(t *testing.T) {
		claims := domain.Claims{
			UserID: 1,
			Email:  "x@example.com",
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			},
		}

		tokenStr, _ := jwtSvc.CreateToken(claims)

		c, err := jwtSvc.ParseToken(tokenStr)
		assert.NoError(t, err)
		assert.Equal(t, uint(1), c.UserID)
		assert.Equal(t, "x@example.com", c.Email)
		assert.True(t, jwtSvc.IsValid(tokenStr))
	})
}
