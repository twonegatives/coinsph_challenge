// Package pgstorage is an PostgreSQL implementation of storage interface
package pgstorage

import (
	"context"
	"database/sql"

	"github.com/shopspring/decimal"
	"github.com/twonegatives/coinsph_challenge/pkg/entities"
)

type PgStorage struct {
	DB *sql.DB
}

func NewPgStorage(db *sql.DB) *PgStorage {
	return &PgStorage{DB: db}
}

func (s *PgStorage) GetAccountsList(ctx context.Context) ([]entities.Account, error) {
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

func (s *PgStorage) GetPaymentsList(ctx context.Context) ([]entities.Payment, error) {
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

func (s *PgStorage) SendPayment(ctx context.Context, from entities.Account, to entities.Account) error {
	return nil
}
