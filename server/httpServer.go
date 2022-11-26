package server

import (
	"github.com/julienschmidt/httprouter"
	"github.com/spf13/viper"
	"net/http"
)

type Server struct {
	httpServer *http.Server
}

func NewServer() *Server {
	return &Server{}
}
func (s *Server) Run(router *httprouter.Router) error {
	s.httpServer = &http.Server{
		Addr:    ":" + viper.GetString("app.port"),
		Handler: router,
	}
	return s.httpServer.ListenAndServe()
}
