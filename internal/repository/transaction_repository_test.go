package repository

import (
	"database/sql"
	appErr "github.com/bhuvi1021/TripleA/internal/errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/bhuvi1021/TripleA/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestTransactionRepository_CreateTransaction(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := NewTransactionRepository(db)

	req := models.CreateTransactionArgs{
		SourceAccountId:      1,
		DestinationAccountId: 2,
		Amount:               100.0,
		CurrencyCode:         "INR",
		Reference:            "TXN-123456",
	}

	tests := []struct {
		name       string
		setupMock  func()
		expectErr  error
		expectResp bool
	}{
		{
			name: "success",
			setupMock: func() {
				mock.ExpectBegin()

				mock.ExpectQuery("SELECT balance FROM accounts").
					WithArgs(req.SourceAccountId).
					WillReturnRows(sqlmock.NewRows([]string{"balance"}).AddRow(200.0))

				mock.ExpectQuery("SELECT balance FROM accounts").
					WithArgs(req.DestinationAccountId).
					WillReturnRows(sqlmock.NewRows([]string{"balance"}).AddRow(50.0))

				mock.ExpectExec("UPDATE accounts SET balance").
					WithArgs(100.0, req.SourceAccountId).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectExec("UPDATE accounts SET balance").
					WithArgs(150.0, req.DestinationAccountId).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectExec("INSERT INTO transactions").
					WithArgs(req.SourceAccountId, req.Amount, req.CurrencyCode, false, req.Reference).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectExec("INSERT INTO transactions").
					WithArgs(req.DestinationAccountId, req.Amount, req.CurrencyCode, true, req.Reference).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
			expectErr:  nil,
			expectResp: true,
		},
		{
			name: "insufficient balance",
			setupMock: func() {
				mock.ExpectBegin()
				mock.ExpectQuery("SELECT balance FROM accounts").
					WithArgs(req.SourceAccountId).
					WillReturnRows(sqlmock.NewRows([]string{"balance"}).AddRow(50.0)) // less than amount
			},
			expectErr:  appErr.ErrInsufficientBalance,
			expectResp: false,
		},
		{
			name: "get source balance fails",
			setupMock: func() {
				mock.ExpectBegin()
				mock.ExpectQuery("SELECT balance FROM accounts").
					WithArgs(req.SourceAccountId).
					WillReturnError(sql.ErrConnDone)
			},
			expectErr:  appErr.ErrTransactionFailed,
			expectResp: false,
		},
		{
			name: "begin fails",
			setupMock: func() {
				mock.ExpectBegin().WillReturnError(sql.ErrConnDone)
			},
			expectErr:  appErr.ErrTransactionFailed,
			expectResp: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			resp, err := repo.CreateTransaction(req)
			assert.Equal(t, tt.expectErr, err)

			if tt.expectResp {
				assert.Equal(t, req.SourceAccountId, resp.SourceAccountId)
				assert.Equal(t, "100.00000", resp.AvailableBalance)
			} else {
				assert.Zero(t, resp.SourceAccountId)
			}
		})
	}
}
