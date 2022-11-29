package auth

import (
	"authService/pkg/db"
	"errors"
	"github.com/go-redis/redis/v9"
	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"
	"log"
	"time"
)

var singingKey = viper.GetString("key.jwt")

type TokenManager struct {
	//TODO REDIS
	db  *db.Client
	red *redis.Client
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
	expiresAccess := jwt.NewNumericDate(time.Now().Add(30 * time.Minute))

	accessClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, &TokenClaims{
		&jwt.RegisteredClaims{
			IssuedAt:  issuedAt,
			ExpiresAt: expiresAccess,
		},
		id,
	})
	//Gen Access token
	accessToken, err := accessClaims.SignedString([]byte(singingKey))
	if err != nil {
		log.Printf("[ERROR] Error: %s\n", err.Error())
		return "", err
	}

	return accessToken, nil
}

func (t *TokenManager) ValidateToken(accessToken string) (int, string, error) {
	//var valid bool
	if accessToken == "" {
		return 0, "", nil
	}

	token, err := jwt.ParseWithClaims(accessToken, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid access-token")
		}

		return []byte(singingKey), nil
	})
	if err != nil {
		log.Print(err.Error())
		return 0, "", err
	}
	claims, ok := token.Claims.(*TokenClaims)
	if !ok {

		return 0, "", errors.New("invalid token") //TODO
	}

	return claims.UserId, accessToken, nil
}
