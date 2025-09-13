package grpc

import (
	"context"
	"errors"

	"github.com/golang-jwt/jwt/v4"
	pbCustomer "github.com/rasteiro11/MCABankAuth/gen/proto/go"
	"github.com/rasteiro11/MCABankAuth/src/auth/service"
	"github.com/rasteiro11/MCABankAuth/src/auth/service/models"
	"github.com/rasteiro11/MCABankAuth/src/user/domain"
	userService "github.com/rasteiro11/MCABankAuth/src/user/service"
	"github.com/rasteiro11/PogCore/pkg/config"
)

type grpcServer struct {
	authService service.AuthService
	userService userService.UserService
}

var (
	ErrInvalidToken     = errors.New("error invalid token")
	ErrSignatureInvalid = errors.New("error invalid signature")
)

type Option func(*grpcServer)

var _ pbCustomer.AuthServiceServer = (*grpcServer)(nil)

func WithUserService(userService userService.UserService) Option {
	return func(s *grpcServer) {
		s.userService = userService
	}
}

func WithAuthService(authService service.AuthService) Option {
	return func(s *grpcServer) {
		s.authService = authService
	}
}

func NewService(opts ...Option) pbCustomer.AuthServiceServer {
	s := &grpcServer{}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func (s *grpcServer) GetUser(ctx context.Context, req *pbCustomer.GetUserRequest) (*pbCustomer.GetUserResponse, error) {
	userEntity := &domain.User{ID: uint(req.Id)}
	userDomain, err := s.userService.FindOne(ctx, userEntity)
	if err != nil {
		return nil, err
	}

	return &pbCustomer.GetUserResponse{
		Id:       int32(userDomain.ID),
		Email:    userDomain.Email,
		Document: userDomain.Document,
	}, nil
}

func (s *grpcServer) VerifySession(ctx context.Context, req *pbCustomer.VerifySessionRequest) (*pbCustomer.VerifySessionResponse, error) {
	claims := &models.Claims{}
	token, err := jwt.ParseWithClaims(req.Token, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(config.Instance().RequiredString("JWT_SECRET")), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrSignatureInvalid) {
			return nil, ErrSignatureInvalid
		}
		return nil, err
	}

	if !token.Valid {
		return nil, ErrSignatureInvalid
	}

	return &pbCustomer.VerifySessionResponse{
		UserId: uint64(claims.UserID),
	}, nil
}
