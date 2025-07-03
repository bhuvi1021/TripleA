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
	accountRepo *repository.AccountRepository
}

func NewAccountService(repo *repository.AccountRepository) *AccountService {
	return &AccountService{accountRepo: repo}
}

func (s *AccountService) CreateAccount(ctx context.Context, req models.CreateAccountRequest) error {
	fName := "AccountService.CreateAccount"
	if req.AccountId == 0 {
		fmt.Println("[%s] failed to get account: %w", fName, appErr.ErrAccountNotFound)
		return appErr.ErrAccountNotFound
	}

	initialBalance, err := parseBalance(req.InitialBalance)
	if err != nil {
		fmt.Println("[%s] failed to parse balance: %w", fName, err)
		return appErr.ErrInvalidAmount
	}

	if initialBalance < 0 {
		return appErr.ErrNegativeBalance
	}

	account, err := s.accountRepo.GetByAccountId(req.AccountId)
	fmt.Println(account)
	if err != nil && err != appErr.ErrAccountNotFound {
		fmt.Println("[%s] failed due to: %w", fName, appErr.ErrInternal)
		return appErr.ErrInternal
	}

	if account != nil {
		fmt.Println("[%s] failed due to: %w", fName, appErr.ErrAccountExists)
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
