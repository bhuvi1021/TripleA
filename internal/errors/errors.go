package errors

import (
	"errors"
	"net/http"
)

var (
	ErrInvalidInput                = errors.New("invalid input")
	ErrInvalidAccountId            = errors.New("invalid account id")
	ErrInvalidSourceAccountId      = errors.New("invalid sender account id")
	ErrInvalidDestinationAccountId = errors.New("invalid receiver account id")
	ErrSameSourceAndDestinationId  = errors.New("sender and receiver account ids must be different")
	ErrAccountExists               = errors.New("account already exists")
	ErrAccountNotFound             = errors.New("account not found")
	ErrSourceAccountNotFound       = errors.New("sender account not found")
	ErrDestinationAccountNotFound  = errors.New("receiver account not found")
	ErrNegativeBalance             = errors.New("initial balance cannot be negative")
	ErrInvalidBalanceFmt           = errors.New("invalid initial balance format")
	ErrInvalidJsonFormat           = errors.New("invalid JSON format")
	ErrInvalidAmount               = errors.New("invalid amount")
	ErrAccountCreation             = errors.New("account creation failed")
	ErrFailedToGetBalance          = errors.New("failed to get balance")
	ErrInsufficientBalance         = errors.New("insufficient funds in sender account")
	ErrTransactionFailed           = errors.New("transaction failed due to internal server error. Please contact support team")
	ErrInternal                    = errors.New("internal server error")
)

var HTTPStatusMap = map[error]int{
	ErrInvalidInput:                http.StatusBadRequest,
	ErrInvalidAccountId:            http.StatusBadRequest,
	ErrInvalidSourceAccountId:      http.StatusBadRequest,
	ErrInvalidDestinationAccountId: http.StatusBadRequest,
	ErrSameSourceAndDestinationId:  http.StatusBadRequest,
	ErrInvalidBalanceFmt:           http.StatusBadRequest,
	ErrNegativeBalance:             http.StatusBadRequest,
	ErrAccountExists:               http.StatusConflict,
	ErrAccountNotFound:             http.StatusNotFound,
	ErrSourceAccountNotFound:       http.StatusNotFound,
	ErrDestinationAccountNotFound:  http.StatusNotFound,
	ErrInvalidJsonFormat:           http.StatusBadRequest,
	ErrInvalidAmount:               http.StatusBadRequest,
	ErrInsufficientBalance:         http.StatusBadRequest,
	ErrAccountCreation:             http.StatusInternalServerError,
	ErrFailedToGetBalance:          http.StatusInternalServerError,
	ErrTransactionFailed:           http.StatusInternalServerError,
	ErrInternal:                    http.StatusInternalServerError,
}
