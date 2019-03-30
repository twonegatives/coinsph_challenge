// Package pgstorage is an PostgreSQL implementation of storage interface
package pgstorage

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/twonegatives/coinsph_challenge/pkg/entities"
)

type PgStorage struct {
	DB *sql.DB
}

func NewPgStorage(db *sql.DB) *PgStorage {
	return &PgStorage{DB: db}
}

func (s *PgStorage) CreateAccount(ctx context.Context, accountName string) (entities.Account, error) {
	query := `INSERT INTO accounts(name, balance, currency) VALUES($1, $2, $3) RETURNING id`
	account := entities.Account{
		Name:     accountName,
		Currency: entities.USD,
		Balance:  decimal.New(0, 0),
	}
	err := s.DB.QueryRowContext(ctx, query, account.Name, account.Balance, account.Currency).Scan(&account.ID)
	return account, errors.Wrap(err, "can't create new account")
}

func (s *PgStorage) GetAccountsList(ctx context.Context) ([]entities.Account, error) {
	query := `SELECT id, name, balance, currency FROM accounts`
	rows, err := s.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, errors.Wrap(err, "can't query Accounts list")
	}

	defer rows.Close()

	var accounts []entities.Account
	for rows.Next() {
		var account entities.Account
		err := rows.Scan(&account.ID, &account.Name, &account.Balance, &account.Currency)
		if err != nil {
			return accounts, errors.Wrap(err, "can't scan Account db row")
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
      counterparties.id,
      counterparties.name,
      transactions.id,
      transactions.created_at,
      direction,
      amount
		FROM payments
		INNER JOIN accounts AS owners ON payments.account_id = owners.id
		INNER JOIN accounts AS counterparties ON payments.counterparty_id = counterparties.id
		INNER JOIN transactions ON payments.transaction_id = transactions.id
	`
	rows, err := s.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, errors.Wrap(err, "can't query Payments list")
	}

	defer rows.Close()

	var payments []entities.Payment
	for rows.Next() {
		var payment entities.Payment
		err := rows.Scan(
			&payment.Account.ID,
			&payment.Account.Name,
			&payment.Counterparty.ID,
			&payment.Counterparty.Name,
			&payment.Transaction.ID,
			&payment.Transaction.CreatedAt,
			&payment.Direction,
			&payment.Amount,
		)
		if err != nil {
			return payments, errors.Wrap(err, "can't scan Payment db row")
		}
		payments = append(payments, payment)
	}

	return payments, nil
}

func (s *PgStorage) SendPayment(ctx context.Context, from entities.Account, to entities.Account, amount decimal.Decimal) error {
	if from.Name == to.Name {
		return errors.New("can't transfer funds to the same account")
	}

	tx, err := s.DB.Begin()
	defer tx.Rollback()

	if err != nil {
		return errors.Wrap(err, "can't open transaction")
	}

	selectQuery := "SELECT id FROM accounts WHERE name = $1 FOR UPDATE"
	paymentSides := getSortedPaymentSides(&from, &to)
	for _, side := range paymentSides {
		if err := tx.QueryRowContext(ctx, selectQuery, side.account.Name).Scan(&side.account.ID); err != nil {
			return errors.Wrapf(err, "can't obtain %s account", side.label)
		}
	}

	var txID int
	insertTxQuery := "INSERT INTO transactions(created_at) VALUES(NOW()) RETURNING id"
	if err := tx.QueryRowContext(ctx, insertTxQuery).Scan(&txID); err != nil {
		return errors.Wrap(err, "can't insert new transaction")
	}

	insertPaymentQuery := `
		INSERT INTO payments(
			transaction_id,
			account_id,
			counterparty_id,
			direction,
			amount
		) VALUES ($1, $2, $3, $4, $5)
		`
	if _, err := tx.ExecContext(ctx, insertPaymentQuery, txID, from.ID, to.ID, "outgoing", amount); err != nil {
		return errors.Wrap(err, "can't insert outgoing payment")
	}

	if _, err := tx.ExecContext(ctx, insertPaymentQuery, txID, to.ID, from.ID, "incoming", amount); err != nil {
		return errors.Wrap(err, "can't insert incoming payment")
	}

	if _, err := tx.ExecContext(ctx, "UPDATE accounts SET balance = balance - $1 WHERE id = $2", amount, from.ID); err != nil {
		return errors.Wrap(err, "can't update sender balance")
	}

	if _, err := tx.ExecContext(ctx, "UPDATE accounts SET balance = balance + $1 WHERE id = $2", amount, to.ID); err != nil {
		return errors.Wrap(err, "can't update receiver balance")
	}

	return errors.Wrap(tx.Commit(), "transaction commit failed")
}

type paymentSide struct {
	account *entities.Account
	label   string
}

func getSortedPaymentSides(from *entities.Account, to *entities.Account) []paymentSide {
	sender := paymentSide{account: from, label: "sender"}
	receiver := paymentSide{account: to, label: "receiver"}

	if from.Name > to.Name {
		return []paymentSide{sender, receiver}
	}

	return []paymentSide{receiver, sender}
}
