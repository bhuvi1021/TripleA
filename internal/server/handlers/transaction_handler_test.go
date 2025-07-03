package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	appErr "github.com/bhuvi1021/TripleA/internal/errors"
	"github.com/bhuvi1021/TripleA/internal/models"
	"github.com/bhuvi1021/TripleA/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTransactionHandler_CreateTransaction(t *testing.T) {
	mockSvc := new(mocks.ITransactionService)
	handler := NewTransactionHandler(mockSvc)

	tests := []struct {
		name           string
		requestBody    string
		mockSetup      func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name:        "Success",
			requestBody: `{"source_account_id":1,"destination_account_id":2,"amount":"100.00"}`,
			mockSetup: func() {
				mockSvc.On("CreateTransaction", mock.Anything, models.CreateTransactionArgs{
					SourceAccountId:      1,
					DestinationAccountId: 2,
					Amount:               100.00,
				}).Return(models.CreateTransactionResponse{
					SourceAccountId:  1,
					AvailableBalance: "900.00000",
				}, nil).Once()
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"source_account_id":1,"available_balance":"900.00000"}`,
		},
		{
			name:           "Invalid JSON",
			requestBody:    `{"source_account_id":1,`,
			mockSetup:      func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error_message":"invalid JSON format"}`,
		},
		{
			name:           "Invalid Amount Format",
			requestBody:    `{"source_account_id":1,"destination_account_id":2,"amount":"abc"}`,
			mockSetup:      func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error_message":"invalid amount"}`,
		},
		{
			name:        "Insufficient Balance",
			requestBody: `{"source_account_id":1,"destination_account_id":2,"amount":"100.00"}`,
			mockSetup: func() {
				mockSvc.On("CreateTransaction", mock.Anything, mock.Anything).
					Return(models.CreateTransactionResponse{}, appErr.ErrInsufficientBalance).Once()
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error_message":"insufficient funds in sender account"}`,
		},
		{
			name:        "Internal Error",
			requestBody: `{"source_account_id":1,"destination_account_id":2,"amount":"100.00"}`,
			mockSetup: func() {
				mockSvc.On("CreateTransaction", mock.Anything, mock.Anything).
					Return(models.CreateTransactionResponse{}, appErr.ErrTransactionFailed).Once()
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error_message":"transaction failed due to internal server error. Please contact support team"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			req := httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewBuffer([]byte(tt.requestBody)))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			handler.CreateTransaction(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)
			assert.JSONEq(t, tt.expectedBody, rec.Body.String())

			mockSvc.AssertExpectations(t)
		})
	}
}
