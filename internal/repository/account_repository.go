package repository

import (
	"database/sql"
	"fmt"
	appErr "github.com/bhuvi1021/TripleA/internal/errors"
	"github.com/bhuvi1021/TripleA/internal/models"
)

// AccountRepository handles account database operations
type AccountRepository struct {
	db *sql.DB
}

// NewAccountRepository creates a new account repository
func NewAccountRepository(db *sql.DB) *AccountRepository {
	return &AccountRepository{db: db}
}

// Create creates a new account
func (r *AccountRepository) CreateAccount(account models.Account) error {
	fName := "AccountRepository.CreateAccount"
	query := `INSERT INTO accounts (account_id, balance, created_at, updated_at) VALUES ($1, $2, $3, $4)`
	fmt.Printf("db account %+v. ", account)
	_, err := r.db.Exec(query, account.AccountId, account.Balance, account.CreatedAt, account.UpdatedAt)
	if err != nil {
		fmt.Printf("[%s] failed due to error %+v. ", fName, err)
		return appErr.ErrAccountCreation
	}
	return nil
}

// GetByID retrieves an account by ID
func (r *AccountRepository) GetByAccountId(accountId int64) (*models.Account, error) {
	fName := "AccountRepository.GetByAccountId"
	query := `SELECT id, account_id, balance, created_at, updated_at, deleted_at FROM accounts WHERE account_id = $1`
	row := r.db.QueryRow(query, accountId)

	var account models.Account
	err := row.Scan(&account.Id, &account.AccountId, &account.Balance, &account.CreatedAt, &account.UpdatedAt, &account.DeletedAt)
	if err != nil {
		fmt.Printf("[%s] failed due to error %+v. ", fName, err)
		if err == sql.ErrNoRows {
			return nil, appErr.ErrAccountNotFound
		}
		return nil, appErr.ErrInternal
	}

	return &account, nil
}

// UpdateBalance updates an account's balance within a transaction
func (r *AccountRepository) UpdateBalance(tx *sql.Tx, accountId int64, newBalance float64) error {
	fName := "AccountRepository.UpdateBalance"
	query := `UPDATE accounts SET balance = $1, updated_at = CURRENT_TIMESTAMP WHERE account_id = $2`
	result, err := tx.Exec(query, newBalance, accountId)
	if err != nil {
		fmt.Println("[%s] failed to update balance: %w", fName, err)
		return appErr.ErrInternal
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		fmt.Println("[%s] failed to get rows affected: %w", fName, err)
		return appErr.ErrInternal
	}

	if rowsAffected == 0 {
		fmt.Println("[%s] account not found: %w", fName, err)
		return appErr.ErrAccountNotFound
	}

	return nil
}

// GetBalanceForUpdate retrieves an account's balance with row locking
func (r *AccountRepository) GetBalanceForUpdate(tx *sql.Tx, accountId int64) (float64, error) {
	fName := "AccountRepository.GetBalanceForUpdate"
	query := `SELECT balance FROM accounts WHERE account_id = $1 FOR UPDATE`
	row := tx.QueryRow(query, accountId)

	var balance float64
	err := row.Scan(&balance)
	if err != nil {
		fmt.Printf("[%s] failed due to error %+v. ", fName, err)
		if err == sql.ErrNoRows {
			return 0, appErr.ErrAccountNotFound
		}
		return 0, appErr.ErrInternal
	}

	return balance, nil
}
