package banking

import (
	"context"

	"github.com/shopspring/decimal"
	"github.com/twonegatives/coinsph_challenge/pkg/entities"
)

type BankingService interface {
	GetAccountsList(ctx context.Context, req GetAccountsRequest) ([]entities.Account, error)
	GetPaymentsList(ctx context.Context, req GetPaymentsRequest) ([]entities.Payment, error)
	SendPayment(ctx context.Context, req SendPaymentRequest) (entities.Payment, error)
}

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (svc *Service) GetAccountsList(ctx context.Context, req GetAccountsRequest) ([]entities.Account, error) {
	accounts := []entities.Account{
		{
			ID:       "bob123",
			Balance:  decimal.NewFromFloat(100.0),
			Currency: entities.USD,
		},
		{
			ID:       "alice456",
			Balance:  decimal.NewFromFloat(0.01),
			Currency: entities.USD,
		},
	}

	return accounts, nil
}

func (svc *Service) GetPaymentsList(ctx context.Context, req GetPaymentsRequest) ([]entities.Payment, error) {
	payments := []entities.Payment{
		{
			Account:   "bob123",
			Amount:    decimal.NewFromFloat(100.0),
			ToAccount: "alice456",
			Direction: entities.Outgoing,
		},
		{
			Account:     "alice456",
			Amount:      decimal.NewFromFloat(100.0),
			FromAccount: "bob123",
			Direction:   entities.Incoming,
		},
	}

	return payments, nil
}

func (svc *Service) SendPayment(ctx context.Context, req SendPaymentRequest) (entities.Payment, error) {
	payment := entities.Payment{
		Account:   req.From.ID,
		ToAccount: req.To.ID,
		Direction: entities.Outgoing,
		Amount:    req.Amount,
	}

	return payment, nil
}
