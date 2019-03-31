// Package banking provides functionality for performing common banking operations
// like creating account, transferring money and checking current account/payments.
package banking

import (
	"context"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/twonegatives/coinsph_challenge/pkg/entities"
	"github.com/twonegatives/coinsph_challenge/pkg/storage"
)

//go:generate mockgen -source=service.go -destination ../mocks/mock_banking_service.go -package mocks

// BankingService is an abstraction which contains declarations of methods
// used to create/show Accounts and Payments.
type BankingService interface {
	CreateAccount(ctx context.Context, accountName string) (entities.Account, error)
	GetAccountsList(ctx context.Context) ([]entities.Account, error)
	GetPaymentsList(ctx context.Context) ([]entities.Payment, error)
	SendPayment(ctx context.Context, from entities.Account, to entities.Account, amount decimal.Decimal) error
}

// Service is an implementation of BankingService.
type Service struct {
	store storage.Storage
}

func NewService(s storage.Storage) *Service {
	return &Service{
		store: s,
	}
}

// CreateAccount method accepts a Account object with Name field filled in.
// Tries to create a new account with this name. Returns Account entity with
// all the attributes set up on success.
func (svc *Service) CreateAccount(ctx context.Context, accountName string) (entities.Account, error) {
	account, err := svc.store.CreateAccount(ctx, accountName)
	return account, errors.Wrap(err, "failed to create new account in database")
}

// GetAccountsList returns all the accounts which currently exist in system.
func (svc *Service) GetAccountsList(ctx context.Context) ([]entities.Account, error) {
	accounts, err := svc.store.GetAccountsList(ctx)
	return accounts, errors.Wrap(err, "failed to fetch accounts list from database")
}

// GetPaymentsList returns all the payments which currently exist in system.
func (svc *Service) GetPaymentsList(ctx context.Context) ([]entities.Payment, error) {
	payments, err := svc.store.GetPaymentsList(ctx)
	return payments, errors.Wrap(err, "failed to fetch payments list from database")
}

// SendPayment attempts to transfer 'amount' of money between 'from' and 'to' Accounts.
// Returns error in the following cases:
// - 'from' and 'to' are the same account
// - 'from' has insufficient funds (balance would go < 0 after transfer)
// - either 'from' or 'to' account is not present in system
// - there is an existing mismatch between Account's balance and his/her payments (checked by db trigger)
func (svc *Service) SendPayment(ctx context.Context, from entities.Account, to entities.Account, amount decimal.Decimal) error {
	// TODO: consider wrapping errors to something application-specific
	if from == to {
		return errors.New("can't transfer funds to the same account")
	}

	txStorage, err := svc.store.BeginTx(ctx, nil)
	defer txStorage.RollbackTx(ctx)

	if err != nil {
		return errors.Wrap(err, "can't open transaction")
	}

	// We need to lock Account rows safely in a determined order.
	// By having sender and receiver sorted by name and locked in this
	// order we ensure that our code is not a subject to a deadlock
	paymentSides := getSortedPaymentSides(&from, &to)
	for _, side := range paymentSides {
		if err := txStorage.GetAccountForUpdate(ctx, side.account); err != nil {
			return errors.Wrapf(err, "can't obtain %s account", side.label)
		}
	}

	if from.Balance.LessThan(amount) && !from.MayGoBelowZero() {
		return errors.New("sender account has insufficient funds")
	}

	transaction, err := txStorage.CreateTransaction(ctx)
	if err != nil {
		return errors.Wrap(err, "can't insert new transaction")
	}

	outgoingPayment := entities.Payment{
		Account:      from,
		Counterparty: to,
		Amount:       amount,
		Transaction:  transaction,
		Direction:    entities.Outgoing,
		Currency:     entities.USD,
	}
	if err := txStorage.SendPayment(ctx, outgoingPayment); err != nil {
		return errors.Wrap(err, "can't insert outgoing payment")
	}

	incomingPayment := entities.Payment{
		Account:      to,
		Counterparty: from,
		Amount:       amount,
		Transaction:  transaction,
		Direction:    entities.Incoming,
		Currency:     entities.USD,
	}

	if err := txStorage.SendPayment(ctx, incomingPayment); err != nil {
		return errors.Wrap(err, "can't insert incoming payment")
	}

	from.Balance = from.Balance.Sub(amount)
	if err := txStorage.SetAccountBalance(ctx, from); err != nil {
		return errors.Wrap(err, "can't update sender account balance")
	}

	to.Balance = to.Balance.Add(amount)
	if err := txStorage.SetAccountBalance(ctx, to); err != nil {
		return errors.Wrap(err, "can't update counterparty balance")
	}

	return errors.Wrap(txStorage.CommitTx(ctx), "transaction commit failed")
}

type paymentSide struct {
	account *entities.Account
	label   string
}

// sorts [from, to] accounts by name and returns a slice in determined order
func getSortedPaymentSides(from *entities.Account, to *entities.Account) []paymentSide {
	sender := paymentSide{account: from, label: "sender"}
	receiver := paymentSide{account: to, label: "receiver"}

	if from.Name > to.Name {
		return []paymentSide{sender, receiver}
	}

	return []paymentSide{receiver, sender}
}
