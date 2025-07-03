package service

import (
	"context"
	"errors"
	"testing"

	appErr "github.com/bhuvi1021/TripleA/internal/errors"
	"github.com/bhuvi1021/TripleA/internal/models"
	"github.com/bhuvi1021/TripleA/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTransactionService_CreateTransaction(t *testing.T) {
	mockTxnRepo := new(mocks.ITransactionRepository)
	mockAcctRepo := new(mocks.IAccountRepository)
	svc := NewTransactionService(mockTxnRepo, mockAcctRepo)

	ctx := context.Background()

	tests := []struct {
		name                 string
		req                  models.CreateTransactionArgs
		mockAccounts         map[int64]*models.Account
		repoError            error
		expectedError        error
		expectedAvailableBal string
	}{
		{
			name:          "invalid source account ID",
			req:           models.CreateTransactionArgs{SourceAccountId: 0, DestinationAccountId: 2, Amount: 100},
			expectedError: appErr.ErrInvalidSourceAccountId,
		},
		{
			name:          "invalid destination account ID",
			req:           models.CreateTransactionArgs{SourceAccountId: 1, DestinationAccountId: 0, Amount: 100},
			expectedError: appErr.ErrInvalidDestinationAccountId,
		},
		{
			name:          "same source and destination account ID",
			req:           models.CreateTransactionArgs{SourceAccountId: 1, DestinationAccountId: 1, Amount: 100},
			expectedError: appErr.ErrSameSourceAndDestinationId,
		},
		{
			name:          "negative amount",
			req:           models.CreateTransactionArgs{SourceAccountId: 1, DestinationAccountId: 2, Amount: -50},
			expectedError: appErr.ErrInvalidAmount,
		},
		{
			name: "source account not found",
			req:  models.CreateTransactionArgs{SourceAccountId: 1, DestinationAccountId: 2, Amount: 100},
			mockAccounts: map[int64]*models.Account{
				1: nil,
			},
			expectedError: appErr.ErrSourceAccountNotFound,
		},
		{
			name: "destination account not found",
			req:  models.CreateTransactionArgs{SourceAccountId: 1, DestinationAccountId: 2, Amount: 100},
			mockAccounts: map[int64]*models.Account{
				1: {},
				2: nil,
			},
			expectedError: appErr.ErrDestinationAccountNotFound,
		},
		{
			name: "repository error during transaction",
			req:  models.CreateTransactionArgs{SourceAccountId: 1, DestinationAccountId: 2, Amount: 100},
			mockAccounts: map[int64]*models.Account{
				1: {},
				2: {},
			},
			repoError:     errors.New("db error"),
			expectedError: errors.New("db error"),
		},
		{
			name: "successful transaction",
			req:  models.CreateTransactionArgs{SourceAccountId: 1, DestinationAccountId: 2, Amount: 100},
			mockAccounts: map[int64]*models.Account{
				1: {},
				2: {},
			},
			expectedAvailableBal: "900.00000",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockAcctRepo.ExpectedCalls = nil
			mockTxnRepo.ExpectedCalls = nil

			if tc.mockAccounts != nil {
				mockAcctRepo.On("GetByAccountId", tc.req.SourceAccountId).
					Return(tc.mockAccounts[tc.req.SourceAccountId], nil).Once()
				mockAcctRepo.On("GetByAccountId", tc.req.DestinationAccountId).
					Return(tc.mockAccounts[tc.req.DestinationAccountId], nil).Once()
			}

			if tc.repoError != nil {
				mockTxnRepo.On("CreateTransaction", mock.Anything).
					Return(models.CreateTransactionResponse{}, tc.repoError).Once()
			} else if tc.expectedAvailableBal != "" {
				mockTxnRepo.On("CreateTransaction", mock.Anything).
					Return(models.CreateTransactionResponse{
						SourceAccountId:  tc.req.SourceAccountId,
						AvailableBalance: tc.expectedAvailableBal,
					}, nil).Once()
			}

			resp, err := svc.CreateTransaction(ctx, tc.req)

			if tc.expectedError != nil {
				assert.EqualError(t, err, tc.expectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.req.SourceAccountId, resp.SourceAccountId)
				assert.Equal(t, tc.expectedAvailableBal, resp.AvailableBalance)
			}
		})
	}
}
