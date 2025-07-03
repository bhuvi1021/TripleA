package models

import (
	"database/sql"
	"time"
)

// Transaction represents a money transfer transaction
type Transaction struct {
	Id           int64        `db:"id"`
	AccountId    int64        `db:"account_id"`
	Amount       float64      `db:"amount"`
	CurrencyCode string       `db:"currency_code"`
	IsCredit     bool         `db:"is_credit"`
	Reference    string       `db:"reference"`
	CreatedAt    time.Time    `db:"created_at"`
	UpdatedAt    time.Time    `db:"updated_at"`
	DeletedAt    sql.NullTime `db:"deleted_at"`
}

// CreateTransactionRequest represents the request body for creating a transaction
type CreateTransactionRequest struct {
	SourceAccountId      int64  `json:"source_account_id,binding:required"`
	DestinationAccountId int64  `json:"destination_account_id,binding:required"`
	Amount               string `json:"amount,binding:required"`
}

// CreateTransactionArgs represents the internal service payload for creating a transaction
type CreateTransactionArgs struct {
	SourceAccountId      int64
	DestinationAccountId int64
	Amount               float64
	CurrencyCode         string
	Reference            string
}

// CreateTransactionResponse represents the response body for creating a transaction
type CreateTransactionResponse struct {
	SourceAccountId  int64  `json:"source_account_id"`
	AvailableBalance string `json:"available_balance"`
}
