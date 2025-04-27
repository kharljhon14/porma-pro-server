package api

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/kharljhon14/porma-pro-server/internal/db/sqlc"
)

type createSummmaryRequest struct {
	AccountID int64  `json:"account_id" binding:"required,min=1"`
	Summary   string `json:"summary" binding:"required,max=3000"`
}

func (s *Server) createSummaryHandler(ctx *gin.Context) {
	var req createSummmaryRequest

	err := ctx.ShouldBindBodyWithJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	args := db.CreateSummaryParams{
		AccountID: req.AccountID,
		Summary:   req.Summary,
	}

	summary, err := s.store.CreateSummary(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, summary)
}

type getSummaryRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (s *Server) getSummaryHandler(ctx *gin.Context) {
	var req getSummaryRequest

	err := ctx.ShouldBindUri(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	summary, err := s.store.GetSummary(ctx, req.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, summary)
}
