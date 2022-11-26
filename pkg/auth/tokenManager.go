package auth

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"
	"log"
	"time"
)

var singingKey = viper.GetString("key.jwt")

type TokenManager struct {
}

type TokenClaims struct {
	*jwt.RegisteredClaims
	UserId int `json:"user_id"`
}

func NewTokenManager() TokenManager {
	return TokenManager{}
}

func (t *TokenManager) GenerateToken(id int) (string, error) {
	if id == 0 {
		return "", errors.New("error: id = 0; unauthorized")
	}
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, &TokenClaims{
		&jwt.RegisteredClaims{
			IssuedAt:  &jwt.NumericDate{time.Now()},
			ExpiresAt: &jwt.NumericDate{time.Now().Add(time.Minute * 30)},
		},
		id,
	})
	token, err := claims.SignedString([]byte(singingKey))
	if err != nil {
		log.Printf("[ERROR] Error: %s\n", err.Error())
		return "", err
	}
	return token, nil
}

func (t *TokenManager) ValidateToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid access-token")
		}
		return []byte(singingKey), nil
	})
	if err != nil {
		log.Printf("[ERROR] Error: %s\n", err.Error())
		return 0, err
	}
	claims, ok := token.Claims.(*TokenClaims)
	if !ok {
		return 0, errors.New("invalid token claims") //TODO
	}
	return claims.UserId, nil
}
