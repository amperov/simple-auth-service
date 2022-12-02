package auth

import (
	"authService/pkg/db"
	hash "authService/pkg/tooling"
	"context"
	"errors"
	"github.com/go-redis/redis/v9"
	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"
	"log"
	"time"
)

var singingKey = viper.GetString("key.jwt")

type TokenManager struct {
	db  *db.Client
	red *redis.Client
}

func Error(err error) {
	if err != nil {
		log.Println(err.Error())
		return
	}
}

type TokenClaims struct {
	*jwt.RegisteredClaims
	UserId int `json:"user_id"`
}

func NewTokenManager(client *redis.Client, conn *db.Client) TokenManager {
	return TokenManager{red: client, db: conn}
}

func (t *TokenManager) GenerateToken(id int) (string, error) {
	if id == 0 {
		return "", errors.New("error: id = 0; unauthorized")
	}

	issuedAt := jwt.NewNumericDate(time.Now())
	expiresAccess := jwt.NewNumericDate(time.Now().Add(1 * time.Minute))
	expiresRefresh := jwt.NewNumericDate(time.Now().Add(30 * 24 * time.Hour))

	accessClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, &TokenClaims{
		&jwt.RegisteredClaims{
			IssuedAt:  issuedAt,
			ExpiresAt: expiresAccess,
		},
		id,
	})
	//Gen Access Token
	accessToken, err := accessClaims.SignedString([]byte(singingKey))
	Error(err)

	refreshClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, &TokenClaims{
		&jwt.RegisteredClaims{
			IssuedAt:  issuedAt,
			ExpiresAt: expiresRefresh,
		},
		id,
	})
	//Gen Refresh Token
	refreshToken, err := refreshClaims.SignedString([]byte(singingKey))
	Error(err)

	//Now We get Hash into 20 symbols
	accessCode := hash.Hash(accessToken, refreshToken)
	t.red.Set(context.Background(), accessCode, accessToken, 30*time.Minute)
	t.db.SetRefresh(accessCode, refreshToken, context.Background())

	return accessCode, nil
}

func (t *TokenManager) ValidateToken(accessCode string) (int, string, error) {
	//var valid bool
	var validAccess bool

	if accessCode == "" {
		return 0, "", nil
	}
	//Getting AccessToken by AccessCode
	accessToken := t.red.Get(context.Background(), accessCode)

	//Our access token
	result := accessToken.Val()

	aToken, err := jwt.ParseWithClaims(result, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid access-token")
		}

		return []byte(singingKey), nil
	})
	validAccess = true
	if err != nil {
		validAccess = false
		log.Print(err.Error())
	}
	if validAccess != false {
		claims, ok := aToken.Claims.(*TokenClaims)
		if !ok {
			return 0, "", errors.New("invalid token") //TODO
		}
		return claims.UserId, accessCode, nil
	}
	if validAccess == false {
		refresh, err := t.db.GetRefresh(accessCode, context.Background())
		if err != nil {
			return 0, "", err
		}
		rToken, err := jwt.ParseWithClaims(refresh, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("invalid access-token")
			}
			return []byte(singingKey), nil
		})
		if err != nil {
			return 0, "", err
		}
		claims, ok := rToken.Claims.(*TokenClaims)
		if !ok {
			return 0, "", errors.New("invalid token") //TODO
		}
		return claims.UserId, accessCode, nil
	}

	return 0, "", err
}
