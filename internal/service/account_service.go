package service

import (
	"context"
	"fmt"
	appErr "github.com/bhuvi1021/TripleA/internal/errors"
	"github.com/bhuvi1021/TripleA/internal/models"
	"github.com/bhuvi1021/TripleA/internal/repository"
	"strconv"
	"time"
)

type AccountService struct {
	accountRepo repository.IAccountRepository
}

func NewAccountService(repo repository.IAccountRepository) *AccountService {
	return &AccountService{accountRepo: repo}
}

type IAccountService interface {
	CreateAccount(ctx context.Context, req models.CreateAccountRequest) error
	GetAccount(ctx context.Context, id int64) (*models.Account, error)
}

func (s *AccountService) CreateAccount(ctx context.Context, req models.CreateAccountRequest) error {
	fName := "AccountService.CreateAccount"
	if req.AccountId <= 0 {
		fmt.Printf("[%s] failed to parse the account id: %v", fName, appErr.ErrInvalidAccountId)
		return appErr.ErrInvalidAccountId
	}

	initialBalance, err := parseBalance(req.InitialBalance)
	if err != nil {
		fmt.Printf("[%s] failed to parse balance: %v", fName, err)
		return appErr.ErrInvalidAmount
	}

	if initialBalance < 0 {
		return appErr.ErrNegativeBalance
	}

	account, err := s.accountRepo.GetByAccountId(req.AccountId)
	if err != nil && err != appErr.ErrAccountNotFound {
		fmt.Printf("[%s] failed due to: %v", fName, appErr.ErrInternal)
		return appErr.ErrInternal
	}

	if account != nil {
		fmt.Printf("[%s] failed due to: %v", fName, appErr.ErrAccountExists)
		return appErr.ErrAccountExists
	}

	timeNow := time.Now().UTC()
	return s.accountRepo.CreateAccount(models.Account{
		AccountId: req.AccountId,
		Balance:   initialBalance,
		CreatedAt: timeNow,
		UpdatedAt: timeNow,
	})
}

func (s *AccountService) GetAccount(ctx context.Context, accountID int64) (*models.Account, error) {
	return s.accountRepo.GetByAccountId(accountID)
}

func parseBalance(balanceStr string) (float64, error) {
	return strconv.ParseFloat(balanceStr, 64)
}
