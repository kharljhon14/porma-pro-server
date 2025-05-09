// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package db

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Account struct {
	ID           int64            `json:"id"`
	Email        string           `json:"email"`
	PasswordHash string           `json:"password_hash"`
	FullName     string           `json:"full_name"`
	CreatedAt    pgtype.Timestamp `json:"created_at"`
	UpdatedAt    pgtype.Timestamp `json:"updated_at"`
	IsVerified   bool             `json:"is_verified"`
}

type PersonalInfo struct {
	ID          int64       `json:"id"`
	AccountID   int64       `json:"account_id"`
	FullName    string      `json:"full_name"`
	Email       string      `json:"email"`
	PhoneNumber string      `json:"phone_number"`
	LinkedinUrl pgtype.Text `json:"linkedin_url"`
	PersonalUrl pgtype.Text `json:"personal_url"`
	Country     string      `json:"country"`
	State       string      `json:"state"`
	City        string      `json:"city"`
}

type Summary struct {
	ID        int64  `json:"id"`
	AccountID int64  `json:"account_id"`
	Summary   string `json:"summary"`
}

type WorkExperience struct {
	ID        int64            `json:"id"`
	AccountID int64            `json:"account_id"`
	Role      string           `json:"role"`
	Company   string           `json:"company"`
	Location  string           `json:"location"`
	Summary   string           `json:"summary"`
	StartDate pgtype.Timestamp `json:"start_date"`
	EndDate   pgtype.Timestamp `json:"end_date"`
}
