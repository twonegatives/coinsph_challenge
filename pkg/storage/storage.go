// Package storage contains a set of interfaces to be
// implemented by persistance repositories (e.g. database, memory stores etc)
package storage

import (
	"context"
	"database/sql"

	"github.com/twonegatives/coinsph_challenge/pkg/entities"
)

//go:generate mockgen -source=storage.go -destination ../mocks/mock_storage.go -package mocks

// Storage is an abstraction unifying methods for objects persistance.
type Storage interface {
	BeginTx(ctx context.Context, opts *sql.TxOptions) (Storage, error)
	CommitTx(ctx context.Context) error
	RollbackTx(ctx context.Context) error

	CreateAccount(ctx context.Context, accountName string) (entities.Account, error)
	GetAccountsList(ctx context.Context) ([]entities.Account, error)
	GetPaymentsList(ctx context.Context) ([]entities.Payment, error)

	GetAccountForUpdate(ctx context.Context, account *entities.Account) error
	CreateTransaction(ctx context.Context) (entities.Transaction, error)
	SendPayment(ctx context.Context, payment entities.Payment) error
	SetAccountBalance(ctx context.Context, account entities.Account) error
}

// TransactionBeginner is an abstraction which allows to start db transaction.
type TransactionBeginner interface {
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

// Queryable is an abstraction with common methods to interact with database.
// It is a common interface for both *sql.DB and *sql.Tx.
type Queryable interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row

	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}
