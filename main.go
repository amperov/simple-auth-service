package main

import (
	"authService/internal/business/service"
	"authService/internal/stores"
	"authService/internal/transport/httpv1/user"
	"authService/pkg/auth"
	"authService/pkg/db"
	"authService/server"
	"github.com/julienschmidt/httprouter"
	"github.com/spf13/viper"
	"log"
)

func main() {
	err := InitConfig()
	if err != nil {
		log.Printf("[FATAL] error reading config\n%s\n", err.Error())
		return
	}

	router := httprouter.New()

	server := server.NewServer()

	client, err := db.GetClient()
	if err != nil {
		return
	}

	authStorage := stores.NewAuthStorage(client)

	tm := auth.NewTokenManager()

	authService := service.NewAuthService(authStorage, tm)

	handler := user.NewAuthHandler(authService)

	handler.Register(router)

	server.Run(router)
}

func InitConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	err := viper.ReadInConfig()
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}
