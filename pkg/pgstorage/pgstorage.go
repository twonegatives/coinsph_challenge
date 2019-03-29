// Package pgstorage is an PostgreSQL implementation of storage interface
package pgstorage

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"
	"github.com/twonegatives/coinsph_challenge/pkg/entities"
)

type PgStorage struct {
	DB *sql.DB
}

func NewPgStorage(db *sql.DB) *PgStorage {
	return &PgStorage{DB: db}
}

func (s *PgStorage) GetAccountsList(ctx context.Context) ([]entities.Account, error) {
	query := `SELECT id, name, balance, currency FROM accounts`
	rows, err := s.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, errors.Wrap(err, "Can not query Accounts list")
	}

	defer rows.Close()

	var accounts []entities.Account
	for rows.Next() {
		var account entities.Account
		err := rows.Scan(&account.ID, &account.Name, &account.Balance, &account.Currency)
		if err != nil {
			return accounts, errors.Wrap(err, "Can't scan Account db row")
		}
		accounts = append(accounts, account)
	}

	return accounts, nil
}

func (s *PgStorage) GetPaymentsList(ctx context.Context) ([]entities.Payment, error) {
	query := `
		SELECT
			owners.id,
			owners.name,
			participants.id,
			participants.name,
			transactions.id,
			transactions.created_at,
			direction,
			amount
		FROM payments
		INNER JOIN accounts AS owners ON payments.account_id = owners.id
		INNER JOIN accounts AS participants ON payments.participant_id = participants.id
		INNER JOIN transactions ON payments.transaction_id = transactions.id
	`
	rows, err := s.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, errors.Wrap(err, "Can not query Payments list")
	}

	defer rows.Close()

	var payments []entities.Payment
	for rows.Next() {
		var payment entities.Payment
		err := rows.Scan(
			&payment.Account.ID,
			&payment.Account.Name,
			&payment.Participant.ID,
			&payment.Participant.Name,
			&payment.Transaction.ID,
			&payment.Transaction.CreatedAt,
			&payment.Direction,
			&payment.Amount,
		)
		if err != nil {
			return payments, errors.Wrap(err, "Can't scan Payment db row")
		}
		payments = append(payments, payment)
	}

	return payments, nil
}

func (s *PgStorage) SendPayment(ctx context.Context, from entities.Account, to entities.Account) error {
	return nil
}
