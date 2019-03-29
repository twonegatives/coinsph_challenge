package banking

import (
	"context"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/twonegatives/coinsph_challenge/pkg/entities"
	"github.com/twonegatives/coinsph_challenge/pkg/storage"
)

type BankingService interface {
	GetAccountsList(ctx context.Context, req GetAccountsRequest) ([]entities.Account, error)
	GetPaymentsList(ctx context.Context, req GetPaymentsRequest) ([]entities.Payment, error)
	SendPayment(ctx context.Context, req SendPaymentRequest) error
}

type Service struct {
	store storage.Storage
}

func NewService(s storage.Storage) *Service {
	return &Service{
		store: s,
	}
}

func (svc *Service) GetAccountsList(ctx context.Context, _req GetAccountsRequest) ([]entities.Account, error) {
	accounts, err := svc.store.GetAccountsList(ctx)
	return accounts, errors.Wrap(err, "failed to fetch accounts list from database")
}

func (svc *Service) GetPaymentsList(ctx context.Context, _req GetPaymentsRequest) ([]entities.Payment, error) {
	payments, err := svc.store.GetPaymentsList(ctx)
	return payments, errors.Wrap(err, "failed to fetch payments list from database")
}

func (svc *Service) SendPayment(ctx context.Context, _req SendPaymentRequest) error {
	err := svc.store.SendPayment(ctx, entities.Account{}, entities.Account{}, decimal.New(0, 0))
	return errors.Wrap(err, "failed to save new payment")
}
