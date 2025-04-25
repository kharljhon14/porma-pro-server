package api

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/kharljhon14/porma-pro-server/internal/db/sqlc"
)

type createPersonalInfoRequest struct {
	AccountID   int64  `json:"account_id" binding:"required,min=1"`
	Email       string `json:"email" binding:"required,email"`
	FullName    string `json:"full_name" binding:"required,max=255"`
	PhoneNumber string `json:"phone_number" binding:"required"`
	LinkedInURL string `json:"linkedin_url" binding:"max=255"`
	PersonalURL string `json:"personal_url" binding:"max=255"`
	Country     string `json:"country" binding:"required,max=255"`
	State       string `json:"state" binding:"required,max=255"`
	City        string `json:"city" binding:"required,max=255"`
}

func (s *Server) createPersonalInfoHandler(ctx *gin.Context) {
	var req createPersonalInfoRequest

	err := ctx.ShouldBindBodyWithJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	args := db.CreatePersonalInfoParams{
		AccountID:   req.AccountID,
		Email:       req.Email,
		FullName:    req.FullName,
		PhoneNumber: req.PhoneNumber,
		Country:     req.Country,
		State:       req.State,
		City:        req.City,
	}

	if req.LinkedInURL != "" {
		args.LinkedinUrl = pgtype.Text{String: req.LinkedInURL}
	}

	if req.PersonalURL != "" {
		args.PersonalUrl = pgtype.Text{String: req.PersonalURL}
	}

	personalInfo, err := s.store.CreatePersonalInfo(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, personalInfo)
}

type getPersonalInfoRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (s *Server) getPersonalInfoHandler(ctx *gin.Context) {
	var req getPersonalInfoRequest

	err := ctx.ShouldBindUri(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	personalInfo, err := s.store.GetPersonalInfo(ctx, req.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, personalInfo)
}

type updatePersonalInfoURI struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type updatePersonalInfoRequest struct {
	Email       string `json:"email" binding:"required,email"`
	FullName    string `json:"full_name" binding:"required,max=255"`
	PhoneNumber string `json:"phone_number" binding:"required"`
	LinkedInURL string `json:"linkedin_url" binding:"max=255"`
	PersonalURL string `json:"personal_url" binding:"max=255"`
	Country     string `json:"country" binding:"required,max=255"`
	State       string `json:"state" binding:"required,max=255"`
	City        string `json:"city" binding:"required,max=255"`
}

func (s *Server) updatePersonalInfoHandler(ctx *gin.Context) {
	var req updatePersonalInfoRequest

	err := ctx.ShouldBindBodyWithJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var uri updatePersonalInfoURI
	err = ctx.ShouldBindUri(&uri)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	args := db.UpdatePersonalInfoParams{
		ID:          uri.ID,
		Email:       req.Email,
		FullName:    req.FullName,
		PhoneNumber: req.PhoneNumber,
		Country:     req.Country,
		State:       req.State,
		City:        req.City,
	}

	if req.LinkedInURL != "" {
		args.LinkedinUrl = pgtype.Text{String: req.LinkedInURL}
	}

	if req.PersonalURL != "" {
		args.PersonalUrl = pgtype.Text{String: req.PersonalURL}
	}

	personalInfo, err := s.store.UpdatePersonalInfo(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, personalInfo)
}
