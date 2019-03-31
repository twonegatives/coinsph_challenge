// Package pgstorage is an PostgreSQL implementation of storage interface.
// It allows to persist objects in PostgreSQL database.
package pgstorage

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/twonegatives/coinsph_challenge/pkg/entities"
	"github.com/twonegatives/coinsph_challenge/pkg/storage"
)

// PgStorage is an implementation of Storage interface
type PgStorage struct {
	Handler storage.Queryable
}

func NewPgStorage(handler storage.Queryable) *PgStorage {
	return &PgStorage{Handler: handler}
}

// BeginTx starts db transaction and returns a new Storage implementation with the transaction as a handler for db queries
func (s *PgStorage) BeginTx(ctx context.Context, opts *sql.TxOptions) (storage.Storage, error) {
	dbConn, ok := s.Handler.(storage.TransactionBeginner)
	if !ok {
		return nil, errors.New("handler doesn't satisfy the interface TransactionBeginner")
	}

	tx, err := dbConn.BeginTx(ctx, opts)
	if err != nil {
		return nil, errors.Wrap(err, "failed to start tx")
	}

	newClient := &PgStorage{Handler: tx}
	return newClient, nil
}

// CommitTx commits a current db transaction if exists. Does nothing if transaction was not open.
func (s *PgStorage) CommitTx(ctx context.Context) error {
	tx := s.tx()
	if tx == nil {
		return errors.New("nothing to commit, transaction is not started")
	}

	return errors.Wrap(tx.Commit(), "failed to commit db transaction")
}

// RollbackTx rollbacks a current db transaction if exists. Does nothing if transaction was not open.
func (s *PgStorage) RollbackTx(ctx context.Context) error {
	tx := s.tx()
	if tx == nil {
		return errors.New("nothing to rollback, transaction is not started")
	}

	return errors.Wrap(tx.Rollback(), "failed to rollback db transaction")
}

// tx returns database transaction if the storage started it previously
func (s *PgStorage) tx() *sql.Tx {
	if tx, ok := s.Handler.(*sql.Tx); ok {
		return tx
	}
	return nil
}

// CreateAccount accepts account name as an argument and tries to
// create a new account with such name. Returns the created Account
// object on success.
func (s *PgStorage) CreateAccount(ctx context.Context, accountName string) (entities.Account, error) {
	query := `INSERT INTO accounts(name, balance, currency) VALUES($1, $2, $3) RETURNING id`
	account := entities.Account{
		Name:     accountName,
		Currency: entities.USD,
		Balance:  decimal.New(0, 0),
	}
	err := s.Handler.QueryRowContext(ctx, query, account.Name, account.Balance, account.Currency).Scan(&account.ID)
	return account, errors.Wrap(err, "can't create new account")
}

// GetAccountsList returns slice of Accounts currently existing in the system
func (s *PgStorage) GetAccountsList(ctx context.Context) ([]entities.Account, error) {
	query := `SELECT id, name, balance, currency FROM accounts`
	rows, err := s.Handler.QueryContext(ctx, query)
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

// GetPaymentsList returns slice of Payments currently existing in the system
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
			amount,
			payments.currency
		FROM payments
		INNER JOIN accounts AS owners ON payments.account_id = owners.id
		INNER JOIN accounts AS counterparties ON payments.counterparty_id = counterparties.id
		INNER JOIN transactions ON payments.transaction_id = transactions.id
	`
	rows, err := s.Handler.QueryContext(ctx, query)
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
			&payment.Currency,
		)
		if err != nil {
			return payments, errors.Wrap(err, "can't scan Payment db row")
		}
		payments = append(payments, payment)
	}

	return payments, nil
}

// GetAccountForUpdate returns Account entitiy with an explicit declaration of row lock
func (s *PgStorage) GetAccountForUpdate(ctx context.Context, account *entities.Account) error {
	selectQuery := "SELECT id, balance FROM accounts WHERE name = $1 FOR UPDATE"
	err := s.Handler.QueryRowContext(ctx, selectQuery, account.Name).Scan(&account.ID, &account.Balance)
	return errors.Wrapf(err, "can't obtain account %s", account.Name)
}

// CreateTransaction creates a Transaction entity
func (s *PgStorage) CreateTransaction(ctx context.Context) (entities.Transaction, error) {
	var result entities.Transaction
	insertTxQuery := "INSERT INTO transactions(created_at) VALUES(NOW()) RETURNING id"
	err := s.Handler.QueryRowContext(ctx, insertTxQuery).Scan(&result.ID)
	return result, errors.Wrap(err, "can't insert new transaction")
}

// SendPayment creates a single Payment entity.
// It is expected to be called two times per each Transaction.
func (s *PgStorage) SendPayment(ctx context.Context, payment entities.Payment) error {
	insertPaymentQuery := `
		INSERT INTO payments(
			transaction_id,
			account_id,
			counterparty_id,
			direction,
			amount,
			currency
		) VALUES ($1, $2, $3, $4, $5, $6)
		`
	_, err := s.Handler.ExecContext(ctx, insertPaymentQuery, payment.Transaction.ID, payment.Account.ID, payment.Counterparty.ID, payment.Direction, payment.Amount, payment.Currency)
	return errors.Wrapf(err, "can't insert %s payment", payment.Direction)
}

// SetAccountBalance takes a single Account entity and updates the related
// database row with balance equal to incoming Account entity's balance
func (s *PgStorage) SetAccountBalance(ctx context.Context, account entities.Account) error {
	_, err := s.Handler.ExecContext(ctx, "UPDATE accounts SET balance = $1 WHERE id = $2", account.Balance, account.ID)
	return errors.Wrapf(err, "can't update balance of %s", account.Name)
}
