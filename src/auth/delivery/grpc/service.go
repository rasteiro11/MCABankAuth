package grpc

import (
	"context"
	"errors"
	pbCustomer "github.com/rasteiro11/MCABankAuth/gen/proto/go"
	"github.com/rasteiro11/MCABankAuth/src/auth/service"
	"github.com/rasteiro11/MCABankAuth/src/user/domain"
	userService "github.com/rasteiro11/MCABankAuth/src/user/service"
)

type grpcServer struct {
	authService service.AuthService
	userService userService.UserService
}

var (
	ErrInvalidToken = errors.New("error invalid token")
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
		Id:    int32(userDomain.ID),
		Email: userDomain.Email,
	}, nil
}

func (s *grpcServer) VerifySession(ctx context.Context, req *pbCustomer.VerifySessionRequest) (*pbCustomer.VerifySessionResponse, error) {
	token, err := s.authService.VerifyToken(ctx, req.Token)
	if err != nil {
		return nil, err
	}

	return &pbCustomer.VerifySessionResponse{
		UserId: uint64(token.UserID),
	}, nil
}
