package repository

import (
	"database/sql"
	"fmt"
	appErr "github.com/bhuvi1021/TripleA/internal/errors"
	"github.com/bhuvi1021/TripleA/internal/models"
	"strconv"
)

// TransactionRepository handles transaction database operations
type TransactionRepository struct {
	db *sql.DB
}

// NewTransactionRepository creates a new transaction repository
func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

type ITransactionRepository interface {
	CreateTransaction(req models.CreateTransactionArgs) (models.CreateTransactionResponse, error)
}

// Create creates a new transaction and updates account balances atomically
func (r *TransactionRepository) CreateTransaction(req models.CreateTransactionArgs) (models.CreateTransactionResponse, error) {
	// Start transaction
	fName := "TransactionRepository.CreateTransaction"
	var resp models.CreateTransactionResponse

	tx, err := r.db.Begin()
	if err != nil {
		fmt.Printf("[%s] failed due to %v", fName, err)
		return resp, appErr.ErrTransactionFailed
	}
	defer tx.Rollback()

	// Get source account balance with row lock
	sourceBalance, err := r.getAccountRepository().GetBalanceForUpdate(tx, req.SourceAccountId)
	if err != nil {
		fmt.Printf("[%s] failed to get source account balance: %v", fName, err)
		return resp, appErr.ErrTransactionFailed
	}

	// Check if source account has sufficient funds
	if sourceBalance < req.Amount {
		return resp, appErr.ErrInsufficientBalance
	}

	// Get destination account balance with row lock
	destinationBalance, err := r.getAccountRepository().GetBalanceForUpdate(tx, req.DestinationAccountId)
	if err != nil {
		fmt.Printf("[%s] failed to get destination account balance: %v", fName, err)
		return resp, appErr.ErrTransactionFailed
	}

	// Calculate new balances
	newSourceBalance := sourceBalance - req.Amount
	newDestinationBalance := destinationBalance + req.Amount

	// Update source account balance
	if err := r.getAccountRepository().UpdateBalance(tx, req.SourceAccountId, newSourceBalance); err != nil {
		fmt.Printf("[%s] failed to update source account balance: %v", fName, err)
		return resp, appErr.ErrTransactionFailed
	}

	// Update destination account balance
	if err := r.getAccountRepository().UpdateBalance(tx, req.DestinationAccountId, newDestinationBalance); err != nil {
		fmt.Printf("[%s] failed to update destination account balance: %v", fName, err)
		return resp, appErr.ErrTransactionFailed
	}

	// Create debit transaction record
	query := `INSERT INTO transactions (account_id, amount, currency_code,is_credit,reference) VALUES ($1, $2, $3, $4, $5)`
	_, err = tx.Exec(query, req.SourceAccountId, req.Amount, req.CurrencyCode, false, req.Reference)
	if err != nil {
		fmt.Printf("[%s] failed to create debit transaction record: %v", fName, err)
		return resp, appErr.ErrTransactionFailed
	}

	// Create credit transaction record
	_, err = tx.Exec(query, req.DestinationAccountId, req.Amount, req.CurrencyCode, true, req.Reference)
	if err != nil {
		fmt.Printf("[%s] failed to create credit transaction record: %v", fName, err)
		return resp, appErr.ErrTransactionFailed
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		fmt.Printf("[%s] failed to commit transaction: %v", fName, err)
		return resp, appErr.ErrTransactionFailed
	}

	resp.AvailableBalance = strconv.FormatFloat(newSourceBalance, 'f', 5, 64)
	resp.SourceAccountId = req.SourceAccountId
	return resp, nil
}

// getAccountRepository returns an account repository instance
// This is a helper method to access account operations within transactions
func (r *TransactionRepository) getAccountRepository() *AccountRepository {
	return &AccountRepository{db: r.db}
}
