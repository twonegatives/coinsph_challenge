package pgstorage_test

import (
	"time"

	"github.com/shopspring/decimal"
	"github.com/twonegatives/coinsph_challenge/pkg/entities"
	"github.com/twonegatives/coinsph_challenge/pkg/storage"
)

const (
	insertAccountQuery = `
		INSERT INTO accounts(
			name,
			balance,
			currency
		) VALUES($1, $2, $3)
		RETURNING id
	`
	insertPaymentQuery = `
		INSERT INTO payments(
			transaction_id,
			account_id,
			counterparty_id,
			direction,
			amount,
			currency
		) VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	insertTransactionQuery = `
		INSERT INTO transactions(
			created_at
		) VALUES ($1)
		RETURNING id
	`

	selectPaymentsQuery = `
		SELECT
			account_id,
			counterparty_id,
			direction,
			amount
		FROM payments
		WHERE account_id = $1 AND counterparty_id = $2
		OR account_id = $2 AND counterparty_id = $1
		ORDER BY id
	`

	selectAccountBalanceQuery = `
		SELECT balance
		FROM accounts
		WHERE id = $1
	`

	selectTransactionCountQuery = `
		SELECT COUNT(*)
		FROM transactions
		WHERE id = $1
	`

	selectAccountCountQuery = `
		SELECT COUNT(*)
		FROM accounts
		WHERE name != 'SYSTEM'
	`
)

func selectTransactionsCount(db storage.Queryable, txID int) (int, error) {
	var count int
	err := db.QueryRow(selectTransactionCountQuery, txID).Scan(&count)
	return count, err
}

func getPaymentsForParticipants(db storage.Queryable, accountID, participantID int) ([]payment, error) {
	payments := []payment{}
	rows, err := db.Query(selectPaymentsQuery, accountID, participantID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var payment payment
		if err := rows.Scan(&payment.AccountID, &payment.CounterpartyID, &payment.Direction, &payment.Amount); err != nil {
			return nil, err
		}
		payments = append(payments, payment)
	}

	return payments, nil
}

func getUserAccountsCount(db storage.Queryable) (int, error) {
	var count int
	err := db.QueryRow(selectAccountCountQuery).Scan(&count)
	return count, err
}

func getAccountBalance(db storage.Queryable, accountID int) (decimal.Decimal, error) {
	var balance decimal.Decimal
	err := db.QueryRow(selectAccountBalanceQuery, accountID).Scan(&balance)
	return balance, err
}

func createAccount(db storage.Queryable, name string, balance decimal.Decimal) (entities.Account, error) {
	account := entities.Account{
		Name:     name,
		Balance:  balance,
		Currency: entities.USD,
	}
	err := db.QueryRow(insertAccountQuery, account.Name, account.Balance, account.Currency).Scan(&account.ID)
	return account, err
}

func createTransaction(db storage.Queryable) (entities.Transaction, error) {
	transaction := entities.Transaction{
		CreatedAt: time.Now(),
	}

	err := db.QueryRow(insertTransactionQuery, transaction.CreatedAt).Scan(&transaction.ID)
	return transaction, err
}

func createPayment(db storage.Queryable, txID, accID, cptyID int, amount decimal.Decimal) (entities.Payment, error) {
	payment := entities.Payment{
		Transaction:  entities.Transaction{ID: txID},
		Account:      entities.Account{ID: accID},
		Counterparty: entities.Account{ID: cptyID},
		Direction:    entities.Outgoing,
		Amount:       amount,
	}

	_, err := db.Exec(insertPaymentQuery, txID, accID, cptyID, entities.Outgoing, amount, entities.USD)
	return payment, err
}
