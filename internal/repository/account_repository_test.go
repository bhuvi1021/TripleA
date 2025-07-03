package repository

import (
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	appErr "github.com/bhuvi1021/TripleA/internal/errors"
	"testing"
	"time"

	"github.com/bhuvi1021/TripleA/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestAccountRepository_CreateAccount(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := NewAccountRepository(db)

	account := models.Account{
		AccountId: 101,
		Balance:   1000.50,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tests := []struct {
		name        string
		mockQuery   func()
		expectedErr error
	}{
		{
			name: "success",
			mockQuery: func() {
				mock.ExpectExec("INSERT INTO accounts").
					WithArgs(account.AccountId, account.Balance, account.CreatedAt, account.UpdatedAt).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedErr: nil,
		},
		{
			name: "db error",
			mockQuery: func() {
				mock.ExpectExec("INSERT INTO accounts").
					WithArgs(account.AccountId, account.Balance, account.CreatedAt, account.UpdatedAt).
					WillReturnError(sql.ErrConnDone)
			},
			expectedErr: appErr.ErrAccountCreation,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockQuery()
			err := repo.CreateAccount(account)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestAccountRepository_GetByAccountId(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := NewAccountRepository(db)
	accountId := int64(101)

	tests := []struct {
		name        string
		mockQuery   func()
		expectedErr error
		expectNil   bool
	}{
		{
			name: "success",
			mockQuery: func() {
				mock.ExpectQuery("SELECT id, account_id, balance, created_at, updated_at, deleted_at FROM accounts").
					WithArgs(accountId).
					WillReturnRows(sqlmock.NewRows([]string{
						"id", "account_id", "balance", "created_at", "updated_at", "deleted_at",
					}).AddRow(1, accountId, 500.0, time.Now(), time.Now(), sql.NullTime{}))
			},
			expectedErr: nil,
			expectNil:   false,
		},
		{
			name: "not found",
			mockQuery: func() {
				mock.ExpectQuery("SELECT id, account_id, balance").
					WithArgs(accountId).
					WillReturnError(sql.ErrNoRows)
			},
			expectedErr: appErr.ErrAccountNotFound,
			expectNil:   true,
		},
		{
			name: "internal error",
			mockQuery: func() {
				mock.ExpectQuery("SELECT id, account_id, balance").
					WithArgs(accountId).
					WillReturnError(sql.ErrConnDone)
			},
			expectedErr: appErr.ErrInternal,
			expectNil:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockQuery()
			acc, err := repo.GetByAccountId(accountId)
			assert.Equal(t, tt.expectedErr, err)
			if tt.expectNil {
				assert.Nil(t, acc)
			} else {
				assert.NotNil(t, acc)
			}
		})
	}
}

func TestUpdateBalance(t *testing.T) {
	type testCase struct {
		name        string
		setupMock   func(sqlmock.Sqlmock)
		accountId   int64
		newBalance  float64
		expectedErr error
	}

	tests := []testCase{
		{
			name: "successful update",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE accounts").
					WithArgs(100.50, int64(1)).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
			},
			accountId:   1,
			newBalance:  100.50,
			expectedErr: nil,
		},
		{
			name: "no rows affected",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE accounts").
					WithArgs(200.75, int64(2)).
					WillReturnResult(sqlmock.NewResult(0, 0))
				mock.ExpectRollback()
			},
			accountId:   2,
			newBalance:  200.75,
			expectedErr: appErr.ErrAccountNotFound,
		},
		{
			name: "query execution error",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE accounts").
					WithArgs(300.00, int64(3)).
					WillReturnError(errors.New("query error"))
				mock.ExpectRollback()
			},
			accountId:   3,
			newBalance:  300.00,
			expectedErr: appErr.ErrInternal,
		},
		{
			name: "rows affected error",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE accounts").
					WithArgs(400.00, int64(4)).
					WillReturnResult(sqlmock.NewErrorResult(errors.New("rows affected error")))
				mock.ExpectRollback()
			},
			accountId:   4,
			newBalance:  400.00,
			expectedErr: appErr.ErrInternal,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			tc.setupMock(mock)
			tx, err := db.Begin()
			assert.NoError(t, err)

			repo := NewAccountRepository(db)
			err = repo.UpdateBalance(tx, tc.accountId, tc.newBalance)

			assert.Equal(t, tc.expectedErr, err)

			if err == nil {
				assert.NoError(t, tx.Commit())
			} else {
				assert.NoError(t, tx.Rollback())
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetBalanceForUpdate(t *testing.T) {
	type testCase struct {
		name        string
		setupMock   func(sqlmock.Sqlmock)
		accountId   int64
		expectedBal float64
		expectedErr error
	}

	tests := []testCase{
		{
			name: "successful balance fetch",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				rows := sqlmock.NewRows([]string{"balance"}).AddRow(150.75)
				mock.ExpectQuery("SELECT balance FROM accounts").
					WithArgs(int64(1)).
					WillReturnRows(rows)
				mock.ExpectCommit()
			},
			accountId:   1,
			expectedBal: 150.75,
			expectedErr: nil,
		},
		{
			name: "account not found",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery("SELECT balance FROM accounts").
					WithArgs(int64(2)).
					WillReturnError(sql.ErrNoRows)
				mock.ExpectRollback()
			},
			accountId:   2,
			expectedBal: 0,
			expectedErr: appErr.ErrAccountNotFound,
		},
		{
			name: "scan error",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				rows := sqlmock.NewRows([]string{"balance"}).
					AddRow(nil)
				mock.ExpectQuery("SELECT balance FROM accounts").
					WithArgs(int64(3)).
					WillReturnRows(rows)
				mock.ExpectRollback()
			},
			accountId:   3,
			expectedBal: 0,
			expectedErr: appErr.ErrInternal,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			tc.setupMock(mock)

			tx, err := db.Begin()
			assert.NoError(t, err)

			repo := NewAccountRepository(db)
			balance, err := repo.GetBalanceForUpdate(tx, tc.accountId)

			assert.Equal(t, tc.expectedBal, balance)
			assert.Equal(t, tc.expectedErr, err)

			if err == nil {
				assert.NoError(t, tx.Commit())
			} else {
				assert.NoError(t, tx.Rollback())
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
