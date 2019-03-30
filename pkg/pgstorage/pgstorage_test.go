package pgstorage_test

import (
	"context"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/twonegatives/coinsph_challenge/pkg/entities"
	"github.com/twonegatives/coinsph_challenge/pkg/pgstorage"
)

var (
	system = entities.Account{
		ID:       1,
		Name:     "SYSTEM",
		Balance:  decimal.New(0, 0),
		Currency: entities.USD,
	}
)

func TestPGStorageGetAccountsList(t *testing.T) {
	db, closeDB := prepareDB(t)
	defer closeDB()
	pg := pgstorage.NewPgStorage(db)
	ctx := context.Background()

	t.Run("returns accounts list", func(t *testing.T) {
		charlie, err := createAccount(pg.DB, "charlie", decimal.New(1534, -2))
		require.NoError(t, err)

		expected := []entities.Account{system, charlie}
		accounts, err := pg.GetAccountsList(ctx)
		require.NoError(t, err)

		assert.Equal(t, expected, accounts)
	})
}

func TestPGStorageGetPaymentsList(t *testing.T) {
	db, closeDB := prepareDB(t)
	defer closeDB()
	pg := pgstorage.NewPgStorage(db)
	ctx := context.Background()

	t.Run("returns payments list", func(t *testing.T) {
		setup := func() (entities.Account, entities.Transaction, entities.Payment) {
			andy, err := createAccount(pg.DB, "andy", decimal.New(0, 0))
			require.NoError(t, err)

			transaction, err := createTransaction(pg.DB)
			require.NoError(t, err)

			paymentRecord, err := createPayment(pg.DB, transaction.ID, system.ID, andy.ID, decimal.New(8732, -2))
			require.NoError(t, err)

			return andy, transaction, paymentRecord
		}

		andy, transaction, paymentRecord := setup()

		payments, err := pg.GetPaymentsList(ctx)
		require.NoError(t, err)

		assert.Len(t, payments, 1)
		assert.Equal(t, transaction.ID, payments[0].Transaction.ID)
		assert.Equal(t, system.ID, payments[0].Account.ID)
		assert.Equal(t, andy.ID, payments[0].Counterparty.ID)
		assert.Equal(t, paymentRecord.Amount, payments[0].Amount)
		assert.Equal(t, entities.Outgoing, payments[0].Direction)
	})
}

type payment struct {
	AccountID      int
	CounterpartyID int
	Direction      string
	Amount         decimal.Decimal
}

func TestPGStorageSendPayment(t *testing.T) {
	db, closeDB := prepareDB(t)
	defer closeDB()
	pg := pgstorage.NewPgStorage(db)
	ctx := context.Background()

	t.Run("denies transfer to the same account", func(t *testing.T) {
		err := pg.SendPayment(ctx, system, system, decimal.New(12, 0))
		require.Error(t, err)
		assert.Contains(t, err.Error(), "can't transfer funds to the same account")
	})

	t.Run("denies transfer to unexisting account", func(t *testing.T) {
		err := pg.SendPayment(ctx, system, entities.Account{Name: "unexisting"}, decimal.New(12, 0))
		require.Error(t, err)
		assert.Contains(t, err.Error(), "can't obtain receiver account")
	})

	t.Run("denies transfer from unexisting account", func(t *testing.T) {
		err := pg.SendPayment(ctx, entities.Account{Name: "unexisting"}, system, decimal.New(12, 0))
		require.Error(t, err)
		assert.Contains(t, err.Error(), "can't obtain sender account")
	})

	t.Run("transfer which sets sender balance below zero", func(t *testing.T) {
		t.Run("is allowed for SYSTEM account", func(t *testing.T) {
			john, err := createAccount(pg.DB, "john", decimal.New(0, 0))
			require.NoError(t, err)

			amount := decimal.New(9999, 0)
			err = pg.SendPayment(ctx, system, john, amount)
			require.NoError(t, err)

			payments, err := getPaymentsForParticipants(pg.DB, system.ID, john.ID)
			require.NoError(t, err)

			systemBalance, err := getAccountBalance(pg.DB, system.ID)
			require.NoError(t, err)

			johnBalance, err := getAccountBalance(pg.DB, john.ID)
			require.NoError(t, err)

			expectedPayments := []payment{
				{
					AccountID:      system.ID,
					CounterpartyID: john.ID,
					Direction:      "outgoing",
					Amount:         amount,
				},
				{
					AccountID:      john.ID,
					CounterpartyID: system.ID,
					Direction:      "incoming",
					Amount:         amount,
				},
			}

			assert.Equal(t, expectedPayments, payments)
			assert.Equal(t, amount, johnBalance)
			assert.Equal(t, amount.Mul(decimal.NewFromFloat(-1.0)), systemBalance)
		})

		t.Run("is denied for any other account", func(t *testing.T) {
			setup := func() (entities.Account, entities.Account) {
				mary, err := createAccount(pg.DB, "mary", decimal.New(0, 0))
				require.NoError(t, err)

				wendy, err := createAccount(pg.DB, "wendy", decimal.New(0, 0))
				require.NoError(t, err)

				return mary, wendy
			}

			mary, wendy := setup()
			amount := decimal.New(9999, 0)
			err := pg.SendPayment(ctx, mary, wendy, amount)
			require.Error(t, err)
			assert.Contains(t, err.Error(), "can't update sender balance")
		})
	})
}

func TestPGStorageCreateAccount(t *testing.T) {
	db, closeDB := prepareDB(t)
	defer closeDB()
	pg := pgstorage.NewPgStorage(db)
	ctx := context.Background()

	t.Run("creates account", func(t *testing.T) {
		name := "tony"
		account, err := pg.CreateAccount(ctx, name)
		require.NoError(t, err)

		assert.Equal(t, name, account.Name)
		assert.Equal(t, decimal.New(0, 0), account.Balance)
		assert.Equal(t, entities.USD, account.Currency)
	})

	t.Run("does not create account with duplicated name", func(t *testing.T) {
		name := "SYSTEM"
		_, err := pg.CreateAccount(ctx, name)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "can't create new account")
	})
}
