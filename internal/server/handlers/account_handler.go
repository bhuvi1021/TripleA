package handlers

import (
	"encoding/json"
	appErr "github.com/bhuvi1021/TripleA/internal/errors"
	"github.com/bhuvi1021/TripleA/internal/models"
	"github.com/bhuvi1021/TripleA/internal/service"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

// Handler contains the dependencies for HTTP handlers
type AccountHandler struct {
	service *service.AccountService
}

// New creates a new handler instance
func NewAccountHandler(service *service.AccountService) *AccountHandler {
	return &AccountHandler{service: service}
}

// CreateAccount handles POST /accounts
func (ah *AccountHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	var req models.CreateAccountRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ah.sendErrorResponse(w, appErr.ErrInvalidJsonFormat)
		return
	}

	err := ah.service.CreateAccount(r.Context(), req)
	if err != nil {
		ah.sendErrorResponse(w, err)
		return
	}

	ah.sendSuccessResponse(w)
}

// GetAccount handles GET /accounts/{account_id}
func (ah *AccountHandler) GetAccount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accountIDStr := vars["account_id"]

	accountID, err := strconv.ParseInt(accountIDStr, 10, 64)
	if err != nil {
		ah.sendErrorResponse(w, appErr.ErrInvalidAccountId)
		return
	}

	account, err := ah.service.GetAccount(r.Context(), accountID)
	if err != nil {
		ah.sendErrorResponse(w, appErr.ErrAccountNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.GetAccountResponse{
		AccountId: account.AccountId,
		Balance:   strconv.FormatFloat(account.Balance, 'f', 5, 64),
		IsDeleted: account.DeletedAt.Valid,
	})
}

// sendErrorResponse writes an error response
func (ah *AccountHandler) sendErrorResponse(w http.ResponseWriter, err error) {
	statusCode, ok := appErr.HTTPStatusMap[err]
	if !ok {
		statusCode = http.StatusInternalServerError
		err = appErr.ErrInternal
	}

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(models.ErrorResponse{ErrorMessage: err.Error()})
}

// sendErrorResponse writes an error response
func (ah *AccountHandler) sendSuccessResponse(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(struct{}{})
}
