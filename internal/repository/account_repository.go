package repository

import (
	"database/sql"
	appErr "github.com/bhuvi1021/TripleA/internal/errors"
	"github.com/bhuvi1021/TripleA/internal/models"
	"log"
)

// AccountRepository handles account database operations
type AccountRepository struct {
	db *sql.DB
}

// NewAccountRepository creates a new account repository
func NewAccountRepository(db *sql.DB) *AccountRepository {
	return &AccountRepository{db: db}
}

type IAccountRepository interface {
	CreateAccount(account models.Account) error
	GetByAccountId(accountId int64) (*models.Account, error)
	UpdateBalance(tx *sql.Tx, accountId int64, newBalance float64) error
	GetBalanceForUpdate(tx *sql.Tx, accountId int64) (float64, error)
}

// Create creates a new account
func (r *AccountRepository) CreateAccount(account models.Account) error {
	fName := "AccountRepository.CreateAccount"
	query := `INSERT INTO accounts (account_id, balance, created_at, updated_at) VALUES ($1, $2, $3, $4)`
	_, err := r.db.Exec(query, account.AccountId, account.Balance, account.CreatedAt, account.UpdatedAt)
	if err != nil {
		log.Printf("[%s] failed due to error %+v. ", fName, err)
		return appErr.ErrAccountCreation
	}
	return nil
}

// GetByAccountId retrieves an account by AccountId
func (r *AccountRepository) GetByAccountId(accountId int64) (*models.Account, error) {
	fName := "AccountRepository.GetByAccountId"
	query := `SELECT id, account_id, balance, created_at, updated_at, deleted_at FROM accounts WHERE account_id = $1`
	row := r.db.QueryRow(query, accountId)

	var account models.Account
	err := row.Scan(&account.Id, &account.AccountId, &account.Balance, &account.CreatedAt, &account.UpdatedAt, &account.DeletedAt)
	if err != nil {
		log.Printf("[%s] failed due to error %+v. ", fName, err)
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
		log.Printf("[%s] failed to update balance: %v", fName, err)
		return appErr.ErrInternal
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("[%s] failed to get rows affected: %v", fName, err)
		return appErr.ErrInternal
	}

	if rowsAffected == 0 {
		log.Printf("[%s] account not found: %v", fName, err)
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
		log.Printf("[%s] failed due to error %+v. ", fName, err)
		if err == sql.ErrNoRows {
			return 0, appErr.ErrAccountNotFound
		}
		return 0, appErr.ErrInternal
	}

	return balance, nil
}
