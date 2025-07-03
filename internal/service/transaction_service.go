package service

import (
	"context"
	"fmt"
	appErr "github.com/bhuvi1021/TripleA/internal/errors"
	uuid "github.com/google/uuid"

	"github.com/bhuvi1021/TripleA/internal/models"
	"github.com/bhuvi1021/TripleA/internal/repository"
)

type TransactionService struct {
	transactionRepo *repository.TransactionRepository
	accountRepo     *repository.AccountRepository
}

func NewTransactionService(transactionRepo *repository.TransactionRepository, accountRepo *repository.AccountRepository) *TransactionService {
	return &TransactionService{
		transactionRepo: transactionRepo,
		accountRepo:     accountRepo,
	}
}

func (ts *TransactionService) CreateTransaction(ctx context.Context, req models.CreateTransactionArgs) (resp models.CreateTransactionResponse, err error) {
	if err = ts.validateCreateTransactionRequest(req); err != nil {
		return resp, err
	}

	req.CurrencyCode = "USD"
	req.Reference = generateTransactionRef()
	if resp, err = ts.transactionRepo.CreateTransaction(req); err != nil {
		return resp, err
	}

	return resp, nil
}

func (ts *TransactionService) validateCreateTransactionRequest(req models.CreateTransactionArgs) error {
	// Validate account IDs
	if req.SourceAccountId <= 0 {
		return appErr.ErrInvalidSourceAccountId
	}

	if req.DestinationAccountId <= 0 {
		return appErr.ErrInvalidDestinationAccountId
	}

	if req.SourceAccountId == req.DestinationAccountId {
		return appErr.ErrSameSourceAndDestinationId
	}

	if req.Amount < 0 {
		return appErr.ErrInvalidAmount
	}

	sourceAccount, err := ts.accountRepo.GetByAccountId(req.SourceAccountId)
	if err != nil {
		return appErr.ErrInternal
	}
	if sourceAccount == nil || sourceAccount.DeletedAt.Valid {
		return appErr.ErrSourceAccountNotFound
	}

	destAccount, err := ts.accountRepo.GetByAccountId(req.DestinationAccountId)
	if err != nil {
		return appErr.ErrInternal
	}
	if destAccount == nil || destAccount.DeletedAt.Valid {
		return appErr.ErrDestinationAccountNotFound
	}

	return nil
}

func generateTransactionRef() string {
	txnRef := uuid.New()
	return fmt.Sprintf("TXN-%s", txnRef.String())
}
