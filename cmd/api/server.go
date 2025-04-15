package api

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	db "github.com/kharljhon14/porma-pro-server/internal/db/sqlc"
	"github.com/kharljhon14/porma-pro-server/internal/token"
)

type Server struct {
	store      db.Store
	router     *gin.Engine
	tokenMaker token.Maker
}

func NewServer(store db.Store) (*Server, error) {
	tokenMaker, err := token.NewJWTMaker(os.Getenv("JWTSECRET"))
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker %w", err)
	}

	server := &Server{
		store:      store,
		tokenMaker: tokenMaker,
	}

	server.mountRoutes()

	return server, nil
}

func (s *Server) mountRoutes() {
	router := gin.Default()

	router.GET("/health", s.healthCheckHandler)

	router.POST("/sign-up", s.createAccountHandler)
	router.POST("/login", s.loginAccountHandler)

	router.GET("/accounts/:id", s.getAccountHandler)
	router.POST("/verify/:id", s.verifyAccountHandler)
	s.router = router
}

func (s *Server) Start(address string) error {
	return s.router.Run(address)
}
