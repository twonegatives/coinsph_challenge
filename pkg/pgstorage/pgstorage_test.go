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

func setupDependencies(t *testing.T) (*pgstorage.PgStorage, func(), context.Context) {
	db, closeDB := prepareDB(t)
	pg := pgstorage.NewPgStorage(db)
	ctx := context.Background()

	return pg, closeDB, ctx
}

func TestPGStorageGetAccountsList(t *testing.T) {
	pg, closeDB, ctx := setupDependencies(t)
	defer closeDB()

	t.Run("returns accounts list", func(t *testing.T) {
		charlie, err := createAccount(pg.Handler, "charlie", decimal.New(1534, -2))
		require.NoError(t, err)

		expected := []entities.Account{system, charlie}
		accounts, err := pg.GetAccountsList(ctx)
		require.NoError(t, err)

		assert.Equal(t, expected, accounts)
	})
}

func TestPGStorageGetPaymentsList(t *testing.T) {
	pg, closeDB, ctx := setupDependencies(t)
	defer closeDB()

	t.Run("returns payments list", func(t *testing.T) {
		setup := func() (entities.Account, entities.Transaction, entities.Payment) {
			andy, err := createAccount(pg.Handler, "andy", decimal.New(0, 0))
			require.NoError(t, err)

			transaction, err := createTransaction(pg.Handler)
			require.NoError(t, err)

			paymentRecord, err := createPayment(pg.Handler, transaction.ID, system.ID, andy.ID, decimal.New(8732, -2))
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

func TestPGStorageGetAccountForUpdate(t *testing.T) {
	pg, closeDB, ctx := setupDependencies(t)
	defer closeDB()

	t.Run("sets account ID and balance attributes", func(t *testing.T) {
		emma, err := createAccount(pg.Handler, "emma", decimal.New(22, -1))
		require.NoError(t, err)

		entity := entities.Account{
			Name: "emma",
		}

		err = pg.GetAccountForUpdate(ctx, &entity)
		require.NoError(t, err)

		assert.Equal(t, decimal.New(22, -1), entity.Balance)
		assert.Equal(t, emma.ID, entity.ID)
	})
}

func TestPGStorageCreateTransaction(t *testing.T) {
	pg, closeDB, ctx := setupDependencies(t)
	defer closeDB()

	t.Run("creates new Transaction entity", func(t *testing.T) {
		tx, err := pg.CreateTransaction(ctx)
		require.NoError(t, err)

		count, err := selectTransactionsCount(pg.Handler, tx.ID)
		require.NoError(t, err)

		assert.Equal(t, 1, count)
	})
}

func TestPGStorageSendPayment(t *testing.T) {
	t.Run("inserts payment", func(t *testing.T) {
		pg, closeDB, ctx := setupDependencies(t)
		defer closeDB()

		liam, err := createAccount(pg.Handler, "liam", decimal.New(220, -2))
		require.NoError(t, err)

		tx, err := pg.CreateTransaction(ctx)
		require.NoError(t, err)

		payment := entities.Payment{
			Account:      system,
			Counterparty: liam,
			Amount:       decimal.New(15, -1),
			Direction:    entities.Outgoing,
			Currency:     entities.USD,
			Transaction:  tx,
		}

		err = pg.SendPayment(ctx, payment)
		require.NoError(t, err)

		result, err := getPaymentsForParticipants(pg.Handler, system.ID, liam.ID)
		require.NoError(t, err)

		assert.Equal(t, payment.Account.ID, result[0].AccountID)
		assert.Equal(t, payment.Counterparty.ID, result[0].CounterpartyID)
		assert.Equal(t, string(payment.Direction), result[0].Direction)
		assert.Equal(t, payment.Amount, result[0].Amount)
	})
}

func TestPGStorageSetAccountBalance(t *testing.T) {
	t.Run("updates balance", func(t *testing.T) {
		pg, closeDB, ctx := setupDependencies(t)
		defer closeDB()

		liam, err := createAccount(pg.Handler, "liam", decimal.New(220, -2))
		require.NoError(t, err)

		liam.Balance = decimal.New(0, 0)
		err = pg.SetAccountBalance(ctx, liam)
		require.NoError(t, err)

		liamBalance, err := getAccountBalance(pg.Handler, liam.ID)
		require.NoError(t, err)

		assert.Equal(t, decimal.New(0, 0), liamBalance)
	})

	t.Run("is not allowed to set balance != sum(payments)", func(t *testing.T) {
		pg, closeDB, ctx := setupDependencies(t)
		defer closeDB()

		liam, err := createAccount(pg.Handler, "liam", decimal.New(220, -2))
		require.NoError(t, err)

		liam.Balance = decimal.New(250, 0)
		err = pg.SetAccountBalance(ctx, liam)
		require.Error(t, err)

		assert.Contains(t, err.Error(), "does not correspond to its payments")
	})
}

func TestPGStorageTransactions(t *testing.T) {
	t.Run("does not store anything if tx was rolled back", func(t *testing.T) {
		pg, closeDB, ctx := setupDependencies(t)
		defer closeDB()

		txStore, err := pg.BeginTx(ctx, nil)
		require.NoError(t, err)

		txPgStore, ok := txStore.(*pgstorage.PgStorage)
		require.Equal(t, true, ok)

		_, err = createAccount(txPgStore.Handler, "liam", decimal.New(220, -2))
		require.NoError(t, err)

		err = txPgStore.RollbackTx(ctx)
		require.NoError(t, err)

		count, err := getUserAccountsCount(pg.Handler)
		require.NoError(t, err)

		assert.Equal(t, 0, count)
	})

	t.Run("stores values if tx was commited", func(t *testing.T) {
		pg, closeDB, ctx := setupDependencies(t)
		defer closeDB()

		txStore, err := pg.BeginTx(ctx, nil)
		require.NoError(t, err)

		txPgStore, ok := txStore.(*pgstorage.PgStorage)
		require.Equal(t, true, ok)

		logan, err := createAccount(txPgStore.Handler, "logan", decimal.New(22, -1))
		require.NoError(t, err)

		err = txPgStore.CommitTx(ctx)
		require.NoError(t, err)

		txPgStore.RollbackTx(ctx)

		loganBalance, err := getAccountBalance(pg.Handler, logan.ID)
		require.NoError(t, err)

		assert.Equal(t, decimal.New(22, -1), loganBalance)
	})
}

func TestPGStorageCreateAccount(t *testing.T) {
	t.Run("creates account", func(t *testing.T) {
		pg, closeDB, ctx := setupDependencies(t)
		defer closeDB()

		name := "tony"
		account, err := pg.CreateAccount(ctx, name)
		require.NoError(t, err)

		assert.Equal(t, name, account.Name)
		assert.Equal(t, decimal.New(0, 0), account.Balance)
		assert.Equal(t, entities.USD, account.Currency)

		count, err := getUserAccountsCount(pg.Handler)
		require.NoError(t, err)

		assert.Equal(t, 1, count)
	})

	t.Run("does not create account with duplicated name", func(t *testing.T) {
		pg, closeDB, ctx := setupDependencies(t)
		defer closeDB()

		name := "SYSTEM"
		_, err := pg.CreateAccount(ctx, name)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "can't create new account")

		count, err := getUserAccountsCount(pg.Handler)
		require.NoError(t, err)

		assert.Equal(t, 0, count)
	})
}
