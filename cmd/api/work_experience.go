package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/kharljhon14/porma-pro-server/internal/db/sqlc"
)

type createWorkExperienceRequest struct {
	AccountID int64     `json:"account_id" binding:"required,min=1"`
	Role      string    `json:"role" binding:"required,max=255"`
	Company   string    `json:"company" binding:"required,max=255"`
	Location  string    `json:"location" binding:"required,max=255"`
	Summary   string    `json:"summary" binding:"required,max=255"`
	StartDate time.Time `json:"start_date" binding:"required"`
	EndDate   time.Time `json:"end_date"`
}

func (s *Server) createWorkExperienceHandler(ctx *gin.Context) {
	var req createWorkExperienceRequest

	err := ctx.ShouldBindBodyWithJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	args := db.CreateWorkExperienceParams{
		AccountID: req.AccountID,
		Role:      req.Role,
		Company:   req.Company,
		Location:  req.Location,
		Summary:   req.Summary,
		StartDate: pgtype.Timestamp{
			Time:  req.StartDate,
			Valid: true,
		},
	}

	if !req.EndDate.IsZero() {
		args.EndDate = pgtype.Timestamp{
			Time:  req.EndDate,
			Valid: true,
		}
	}

	workExperience, err := s.store.CreateWorkExperience(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, workExperience)
}

func (s *Server) getWorkExperienceHandler(ctx *gin.Context) {

}

func (s *Server) getWorkExperienceListHandler(ctx *gin.Context) {

}

func (s *Server) updateWorkExperienceHandler(ctx *gin.Context) {

}

func (s *Server) deleteWorkExperienceHandler(ctx *gin.Context) {

}
