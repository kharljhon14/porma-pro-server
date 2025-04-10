package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) healthCheckHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, nil)
}
