package handlers

import (
	"encoding/json"
	appErr "github.com/bhuvi1021/TripleA/internal/errors"
	"net/http"
	"strconv"

	"github.com/bhuvi1021/TripleA/internal/models"
	"github.com/bhuvi1021/TripleA/internal/service"
)

type TransactionHandler struct {
	service *service.TransactionService
}

func NewTransactionHandler(service *service.TransactionService) *TransactionHandler {
	return &TransactionHandler{service: service}
}

func (th *TransactionHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var req models.CreateTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		th.sendErrorResponse(w, appErr.ErrInvalidJsonFormat)
		return
	}

	amount, err := strconv.ParseFloat(req.Amount, 64)
	if err != nil {
		th.sendErrorResponse(w, appErr.ErrInvalidAmount)
		return
	}

	resp, err := th.service.CreateTransaction(r.Context(), models.CreateTransactionArgs{
		SourceAccountId:      req.SourceAccountId,
		DestinationAccountId: req.DestinationAccountId,
		Amount:               amount,
	})
	if err != nil {
		// You can improve this by using a custom error map like the account handler
		th.sendErrorResponse(w, err)
		return
	}

	th.sendSuccessResponse(w, resp)
}

// sendErrorResponse writes an error response
func (th *TransactionHandler) sendErrorResponse(w http.ResponseWriter, err error) {
	statusCode, ok := appErr.HTTPStatusMap[err]
	if !ok {
		statusCode = http.StatusInternalServerError
		err = appErr.ErrInternal
	}

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(models.ErrorResponse{ErrorMessage: err.Error()})
}

// sendErrorResponse writes an error response
func (th *TransactionHandler) sendSuccessResponse(w http.ResponseWriter, resp models.CreateTransactionResponse) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
