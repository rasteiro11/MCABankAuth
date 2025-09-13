package grpc

import (
	pbCustomer "github.com/rasteiro11/MCABankAuth/gen/proto/go"
	"github.com/rasteiro11/MCABankAuth/src/user/domain"
)

func MapUserRequestToDomain(req *pbCustomer.GetUserRequest) *domain.User {
	return &domain.User{
		ID: uint(req.Id),
	}
}

func MapDomainToUserResponse(user *domain.User) *pbCustomer.GetUserResponse {
	return &pbCustomer.GetUserResponse{
		Id:       int32(user.ID),
		Email:    user.Email,
		Document: user.Document,
	}
}

func MapDocumentRequestToDomain(req *pbCustomer.GetUserByDocumentRequest) *domain.User {
	return &domain.User{
		Document: req.Document,
	}
}

func MapDomainToDocumentResponse(user *domain.User) *pbCustomer.GetUserByDocumentResponse {
	return &pbCustomer.GetUserByDocumentResponse{
		Id:       int32(user.ID),
		Email:    user.Email,
		Document: user.Document,
	}
}
