package storage

import (
	"context"

	"github.com/shopspring/decimal"
	"github.com/twonegatives/coinsph_challenge/pkg/entities"
)

type Storage interface {
	GetAccountsList(ctx context.Context) ([]entities.Account, error)
	GetPaymentsList(ctx context.Context) ([]entities.Payment, error)
	SendPayment(ctx context.Context, from entities.Account, to entities.Account, amount decimal.Decimal) error
}
