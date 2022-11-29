package grpc

import (
	"authService/internal/transport/inputs"
	"context"
	"fmt"
	"google.golang.org/protobuf/runtime/protoimpl"
	"log"
)

type AuthService interface {
	SignUp(ctx context.Context, input *inputs.CreateInput) (int, error)
	SignIn(ctx context.Context, input *inputs.AuthInput) (string, error)
	IsAuthed(ctx context.Context, token string) (int, string, error)
}
type GRPCServer struct {
	as AuthService
}

func NewGRPCServer(as AuthService) *GRPCServer {
	return &GRPCServer{as: as}
}

func (s *GRPCServer) SignUp(ctx context.Context, sign *SignUpRequest) (*SignUpResponse, error) {
	var SignUpInput inputs.CreateInput
	var Status string
	//var SignUpResp SignUpResponse

	SignUpInput.Username = sign.GetUsername()
	SignUpInput.Email = sign.GetEmail()
	SignUpInput.Password = sign.GetPassword()

	_, err := s.as.SignUp(ctx, &SignUpInput)
	if err != nil {
		Status = fmt.Sprintf("AS: %s", err.Error())
		return nil, err
	}
	//SignUpResp.Status = "success"
	Status = "Success"
	return &SignUpResponse{
		state:         protoimpl.MessageState{},
		sizeCache:     0,
		unknownFields: nil,
		Status:        Status,
	}, nil
}
func (s *GRPCServer) SignIn(ctx context.Context, sign *SignInRequest) (*SignInResponse, error) {
	var SignInInput inputs.AuthInput
	//var SignUpResp SignUpResponse
	SignInInput.Username = sign.GetUsername()
	SignInInput.Email = sign.GetEmail()
	SignInInput.Password = sign.GetPassword()
	tk, err := s.as.SignIn(ctx, &SignInInput)
	if err != nil {
		return nil, err
	}

	return &SignInResponse{AccessCode: tk}, nil
}
func (s *GRPCServer) IsAuth(ctx context.Context, request *IsAuthRequest) (*IsAuthResponse, error) {
	log.Println("Is Authed Reached!")
	var resp IsAuthResponse
	accessCode := request.GetAccessCode()

	UserID, NewAccess, err := s.as.IsAuthed(ctx, accessCode)
	if err != nil {
		resp.AccessCode = err.Error()
		resp.UserID = 0
		resp.Auth = false
		return &resp, nil
	}
	resp.AccessCode = NewAccess
	resp.UserID = int32(UserID)
	if UserID == 0 {
		resp.AccessCode = err.Error()
		resp.Auth = false
		return &resp, err
	}

	resp.Auth = true

	return &resp, nil
}
func (s *GRPCServer) mustEmbedUnimplementedAuthServer() {

}
