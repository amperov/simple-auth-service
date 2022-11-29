package service

import (
	"authService/internal/transport/grpc"
	"authService/internal/transport/inputs"
	"authService/pkg/auth"
	"context"
	"fmt"
	"log"
)

type AuthStorage interface {
	AuthUser(Username, Email, Password string, ctx context.Context) (int, error)
	CreateUser(Username, Email, Password string, ctx context.Context) (int, error)
}
type authService struct {
	as AuthStorage
	tm auth.TokenManager
}

func NewAuthService(as AuthStorage, tm auth.TokenManager) grpc.AuthService {
	return &authService{as: as, tm: tm}
}

func (s *authService) SignUp(ctx context.Context, input *inputs.CreateInput) (int, error) {
	if input.Username == "" || input.Email == "" || input.Password == "" {

		return 0, fmt.Errorf("[DEBUG] struct must not be empty")
	}
	UserID, err := s.as.CreateUser(input.Username, input.Email, input.Password, ctx)
	if err != nil {
		return 0, err
	}
	return UserID, nil
}

func (s *authService) SignIn(ctx context.Context, input *inputs.AuthInput) (string, error) {
	if (input.Username == "" && input.Email == "") || input.Password == "" {
		return "", fmt.Errorf("[DEBUG] username and email must not be empty at one moment")
	}
	if input.Username == "" {
		log.Printf("[INFO] username is empty")
	} else if input.Email == "" {
		log.Printf("[INFO] email is empty")
	}

	UserID, err := s.as.AuthUser(input.Username, input.Email, input.Password, ctx)
	if err != nil {
		return "", err
	}

	token, err := s.tm.GenerateToken(UserID)
	if err != nil {
		return "", err
	}

	return token, nil

}

func (s *authService) IsAuthed(_ context.Context, AccessCode string) (int, string, error) {
	UserID, AccessCode, err := s.tm.ValidateToken(AccessCode)
	if err != nil {
		log.Println("Error Validating token!")
		return 0, "", err
	}

	return UserID, AccessCode, nil
}
