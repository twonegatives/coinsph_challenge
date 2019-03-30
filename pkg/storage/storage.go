// Package storage contains an interface which should be
// implemented by persistance repositories (e.g. database, memory stores etc)
package storage

import (
	"context"

	"github.com/shopspring/decimal"
	"github.com/twonegatives/coinsph_challenge/pkg/entities"
)

//go:generate mockgen -source=storage.go -destination ../mocks/mock_storage.go -package mocks

// Storage is an abstraction unifying methods for objects persistance.
type Storage interface {
	CreateAccount(ctx context.Context, accountName string) (entities.Account, error)
	GetAccountsList(ctx context.Context) ([]entities.Account, error)
	GetPaymentsList(ctx context.Context) ([]entities.Payment, error)
	SendPayment(ctx context.Context, from entities.Account, to entities.Account, amount decimal.Decimal) error
}
