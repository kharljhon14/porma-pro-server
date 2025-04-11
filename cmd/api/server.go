package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/kharljhon14/porma-pro-server/internal/db/sqlc"
)

type Server struct {
	store  db.Store
	router *gin.Engine
}

func NewServer(store db.Store) (*Server, error) {
	server := &Server{
		store: store,
	}

	server.mountRoutes()

	return server, nil
}

func (s *Server) mountRoutes() {
	router := gin.Default()

	router.GET("/health", s.healthCheckHandler)

	router.POST("/sign-up", s.createAccountHandler)
	router.GET("/accounts/:id", s.getAccountHandler)

	s.router = router
}

func (s *Server) Start(address string) error {
	return s.router.Run(address)
}
