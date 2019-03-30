package banking

import (
	"context"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/twonegatives/coinsph_challenge/pkg/entities"
	"github.com/twonegatives/coinsph_challenge/pkg/storage"
)

//go:generate mockgen -source=service.go -destination ../mocks/mock_banking_service.go -package mocks

type BankingService interface {
	CreateAccount(ctx context.Context, accountName string) (entities.Account, error)
	GetAccountsList(ctx context.Context) ([]entities.Account, error)
	GetPaymentsList(ctx context.Context) ([]entities.Payment, error)
	SendPayment(ctx context.Context, from entities.Account, to entities.Account, amount decimal.Decimal) error
}

type Service struct {
	store storage.Storage
}

func NewService(s storage.Storage) *Service {
	return &Service{
		store: s,
	}
}

func (svc *Service) CreateAccount(ctx context.Context, accountName string) (entities.Account, error) {
	account, err := svc.store.CreateAccount(ctx, accountName)
	return account, errors.Wrap(err, "failed to create new account in database")
}

func (svc *Service) GetAccountsList(ctx context.Context) ([]entities.Account, error) {
	accounts, err := svc.store.GetAccountsList(ctx)
	return accounts, errors.Wrap(err, "failed to fetch accounts list from database")
}

func (svc *Service) GetPaymentsList(ctx context.Context) ([]entities.Payment, error) {
	payments, err := svc.store.GetPaymentsList(ctx)
	return payments, errors.Wrap(err, "failed to fetch payments list from database")
}

func (svc *Service) SendPayment(ctx context.Context, from entities.Account, to entities.Account, amount decimal.Decimal) error {
	if from == to {
		return errors.New("can't transfer funds to the same account")
	}
	err := svc.store.SendPayment(ctx, from, to, amount)
	return errors.Wrap(err, "failed to save new payment")
}
