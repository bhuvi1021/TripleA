package models

import (
	"database/sql"
	"time"
)

// Account represents a financial account of the user. It maintains the current balance
type Account struct {
	Id        int64        `db:"id"`
	AccountId int64        `db:"account_id"`
	Balance   float64      `db:"balance"`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt time.Time    `db:"updated_at"`
	DeletedAt sql.NullTime `db:"deleted_at"`
}

// CreateAccountRequest represents the request body for creating an account
type CreateAccountRequest struct {
	AccountId      int64  `json:"account_id,binding:required"`
	InitialBalance string `json:"initial_balance,binding:required"`
}

// GetAccountResponse represents the response body for creating an account
type GetAccountResponse struct {
	AccountId int64  `json:"account_id"`
	Balance   string `json:"balance"`
	IsDeleted bool   `json:"is_deleted,omitempty"`
}
