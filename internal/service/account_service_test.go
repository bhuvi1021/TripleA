package service

import (
	"context"
	"errors"
	"testing"
	"time"

	appErr "github.com/bhuvi1021/TripleA/internal/errors"
	"github.com/bhuvi1021/TripleA/internal/models"
	"github.com/bhuvi1021/TripleA/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAccountService_CreateAccount(t *testing.T) {
	mockRepo := new(mocks.IAccountRepository)
	service := NewAccountService(mockRepo)

	ctx := context.Background()

	tests := []struct {
		name        string
		req         models.CreateAccountRequest
		setupMocks  func()
		expectedErr error
	}{
		{
			name: "Success",
			req: models.CreateAccountRequest{
				AccountId:      101,
				InitialBalance: "500.00",
			},
			setupMocks: func() {
				mockRepo.On("GetByAccountId", int64(101)).
					Return(nil, appErr.ErrAccountNotFound).Once()
				mockRepo.On("CreateAccount", mock.AnythingOfType("models.Account")).
					Return(nil).Once()
			},
			expectedErr: nil,
		},
		{
			name: "Invalid Account ID",
			req: models.CreateAccountRequest{
				AccountId:      0,
				InitialBalance: "100.00",
			},
			setupMocks:  func() {},
			expectedErr: appErr.ErrAccountNotFound,
		},
		{
			name: "Invalid Balance Format",
			req: models.CreateAccountRequest{
				AccountId:      101,
				InitialBalance: "bad-format",
			},
			setupMocks:  func() {},
			expectedErr: appErr.ErrInvalidAmount,
		},
		{
			name: "Negative Balance",
			req: models.CreateAccountRequest{
				AccountId:      101,
				InitialBalance: "-10",
			},
			setupMocks:  func() {},
			expectedErr: appErr.ErrNegativeBalance,
		},
		{
			name: "Internal DB Error",
			req: models.CreateAccountRequest{
				AccountId:      101,
				InitialBalance: "100",
			},
			setupMocks: func() {
				mockRepo.On("GetByAccountId", int64(101)).
					Return(nil, errors.New("db error")).Once()
			},
			expectedErr: appErr.ErrInternal,
		},
		{
			name: "Account Already Exists",
			req: models.CreateAccountRequest{
				AccountId:      101,
				InitialBalance: "100",
			},
			setupMocks: func() {
				mockRepo.On("GetByAccountId", int64(101)).
					Return(&models.Account{AccountId: 101}, nil).Once()
			},
			expectedErr: appErr.ErrAccountExists,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()
			err := service.CreateAccount(ctx, tt.req)
			assert.Equal(t, tt.expectedErr, err)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestAccountService_GetAccount(t *testing.T) {
	mockRepo := new(mocks.IAccountRepository)
	service := NewAccountService(mockRepo)

	ctx := context.Background()

	account := &models.Account{
		AccountId: 123,
		Balance:   1000.0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockRepo.On("GetByAccountId", int64(123)).Return(account, nil)

	result, err := service.GetAccount(ctx, 123)
	assert.NoError(t, err)
	assert.Equal(t, account, result)
	mockRepo.AssertExpectations(t)
}
