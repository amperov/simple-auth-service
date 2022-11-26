package stores

import (
	"authService/internal/business/service"
	"authService/pkg/db"
	"context"
	"fmt"
	"log"
)

type authStorage struct {
	client *db.Client
}

func NewAuthStorage(client *db.Client) service.AuthStorage {
	return &authStorage{client: client}
}

// CreateUser
//
//	Getting:
//		-Username
//		-Email
//		-Password
//	Returning:
//		-ID from DB
//		-Error or nil
func (s *authStorage) CreateUser(Username, Email, Password string, ctx context.Context) (int, error) {
	exist, err := s.client.IsExist(Username, Email, Password, ctx)
	if err != nil {
		log.Println(err.Error())
		return 0, err
	}

	if exist {

		return 0, fmt.Errorf("[DEBUG] user[%s] exist", Username)
	}

	UserID, err := s.client.Insert(Username, Email, Password, ctx)
	if err != nil {
		log.Println(err.Error())
		return 0, err
	}

	return UserID, nil
}

// AuthUser
//
//	Getting:
//		-Username
//		-Email
//		-Password
//	Returning:
//		-ID from DB
//		-Error or nil
func (s *authStorage) AuthUser(Username, Email, Password string, ctx context.Context) (int, error) {
	//TODO Exist?
	exist, err := s.client.IsExist(Username, Email, Password, ctx)
	if err != nil {
		return 0, err
	}
	if !exist {
		return 0, fmt.Errorf("[DEBUG] user not exist")
	}
	UserID, err := s.client.Get(Username, Email, Password, ctx)
	if err != nil {
		return 0, err
	}
	return UserID, nil
}
func (s *authStorage) DeleteUser() {

}
