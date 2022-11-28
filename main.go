package main

import (
	"authService/internal/business/service"
	"authService/internal/stores"
	grpc2 "authService/internal/transport/grpc"
	"authService/pkg/auth"
	"authService/pkg/db"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	err := InitConfig()
	if err != nil {
		log.Printf("[FATAL] error reading config\n%s\n", err.Error())
		return
	}
	listen, err := net.Listen("tcp", ":8080")
	if err != nil {
		return
	}
	client, err := db.GetClient()
	if err != nil {
		return
	}

	authStorage := stores.NewAuthStorage(client)

	tm := auth.NewTokenManager()

	authService := service.NewAuthService(authStorage, tm)
	newServer := grpc.NewServer()

	grpcServer := grpc2.NewGRPCServer(authService)

	grpc2.RegisterAuthServer(newServer, grpcServer)

	err = newServer.Serve(listen)
	if err != nil {
		log.Fatalln(err.Error())
		return
	}
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
