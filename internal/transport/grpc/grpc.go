package grpc

import (
	"authService/internal/transport/inputs"
	"context"
	"google.golang.org/protobuf/runtime/protoimpl"
)

type AuthService interface {
	SignUp(ctx context.Context, input *inputs.CreateInput) (int, error)
	SignIn(ctx context.Context, input *inputs.AuthInput) (string, error)
	IsAuthed(ctx context.Context, token string) (int, error)
}
type GRPCServer struct {
	as AuthService
}

func NewGRPCServer(as AuthService) *GRPCServer {
	return &GRPCServer{as: as}
}

func (s *GRPCServer) SignUp(ctx context.Context, sign *SignUpRequest) (*SignUpResponse, error) {
	var SignUpInput inputs.CreateInput
	//var SignUpResp SignUpResponse
	SignUpInput.Username = sign.Username
	SignUpInput.Email = sign.Email
	SignUpInput.Password = sign.Password
	_, err := s.as.SignUp(ctx, &SignUpInput)
	if err != nil {
		return nil, err
	}
	//SignUpResp.Status = "success"

	return &SignUpResponse{
		state:         protoimpl.MessageState{},
		sizeCache:     0,
		unknownFields: nil,
		Status:        "Success",
	}, nil
}
func (s *GRPCServer) SignIn(ctx context.Context, sign *SignInRequest) (*SignInResponse, error) {
	var SignInInput inputs.AuthInput
	//var SignUpResp SignUpResponse
	SignInInput.Username = sign.Username
	SignInInput.Email = sign.Email
	SignInInput.Password = sign.Password
	tk, err := s.as.SignIn(ctx, &SignInInput)
	if err != nil {
		return nil, err
	}
	//SignUpResp.Status = "success"

	return &SignInResponse{Token: tk}, nil
}
func (s GRPCServer) mustEmbedUnimplementedAuthServer() {

}
