package handlers

import (
	"bytes"
	"errors"
	"github.com/bhuvi1021/TripleA/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAccountHandler_CreateAccount(t *testing.T) {
	mockService := new(mocks.IAccountService)
	handler := NewAccountHandler(mockService)

	tests := []struct {
		name           string
		inputBody      string
		mockSetup      func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name:      "success",
			inputBody: `{"account_id":101,"initial_balance":"1000"}`,
			mockSetup: func() {
				mockService.On("CreateAccount", mock.Anything, mock.Anything).
					Return(nil).Once()
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{}`,
		},
		{
			name:           "invalid JSON",
			inputBody:      `{"account_id":101,`, // malformed
			mockSetup:      func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error_message":"invalid JSON format"}`,
		},
		{
			name:      "service error - account exists",
			inputBody: `{"account_id":101,"initial_balance":"1000"}`,
			mockSetup: func() {
				mockService.On("CreateAccount", mock.Anything, mock.Anything).
					Return(errors.New("account already exists")).Once()
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error_message":"internal server error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			req := httptest.NewRequest("POST", "/accounts", bytes.NewBuffer([]byte(tt.inputBody)))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			handler.CreateAccount(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.JSONEq(t, tt.expectedBody, rr.Body.String())
			mockService.AssertExpectations(t)
		})
	}
}
