package banking_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/twonegatives/coinsph_challenge/pkg/banking"
	"github.com/twonegatives/coinsph_challenge/pkg/entities"
	"github.com/twonegatives/coinsph_challenge/pkg/mocks"
)

var (
	ErrDB = errors.New("db error")
	ctx   = context.Background()
)

func TestBankingSvcCreateAccount(t *testing.T) {
	t.Run("returns new account", func(t *testing.T) {
		mCtrl := gomock.NewController(t)
		defer mCtrl.Finish()
		storage := mocks.NewMockStorage(mCtrl)

		accName := "bunny"
		storageResult := entities.Account{Name: accName}
		storage.EXPECT().CreateAccount(ctx, accName).Return(storageResult, nil)

		account, err := banking.NewService(storage).CreateAccount(ctx, accName)
		require.NoError(t, err)
		assert.Equal(t, storageResult, account)
	})

	t.Run("propagates storage exceptions", func(t *testing.T) {
		mCtrl := gomock.NewController(t)
		defer mCtrl.Finish()
		storage := mocks.NewMockStorage(mCtrl)

		accName := "duplicated_name"
		storage.EXPECT().CreateAccount(ctx, accName).Return(entities.Account{}, ErrDB)

		_, err := banking.NewService(storage).CreateAccount(ctx, accName)
		require.Error(t, err)
	})
}
func TestBankingSvcGetAccountsList(t *testing.T) {
	t.Run("returns accounts list", func(t *testing.T) {
		mCtrl := gomock.NewController(t)
		defer mCtrl.Finish()
		storage := mocks.NewMockStorage(mCtrl)

		storageResult := []entities.Account{
			{
				Name:     "SYSTEM",
				Balance:  decimal.New(1500, 0),
				Currency: entities.USD,
			},
		}
		storage.EXPECT().GetAccountsList(ctx).Return(storageResult, nil)

		accounts, err := banking.NewService(storage).GetAccountsList(ctx)
		require.NoError(t, err)
		assert.Equal(t, storageResult, accounts)
	})

	t.Run("propagates storage exceptions", func(t *testing.T) {
		mCtrl := gomock.NewController(t)
		defer mCtrl.Finish()
		storage := mocks.NewMockStorage(mCtrl)

		storage.EXPECT().GetAccountsList(ctx).Return(nil, ErrDB)

		_, err := banking.NewService(storage).GetAccountsList(ctx)
		require.Error(t, err)
	})
}
func TestBankingSvcGetPaymentsList(t *testing.T) {
	t.Run("returns payments list", func(t *testing.T) {
		mCtrl := gomock.NewController(t)
		defer mCtrl.Finish()
		storage := mocks.NewMockStorage(mCtrl)

		storageResult := []entities.Payment{
			{
				Account:      entities.Account{Name: "businessman"},
				Counterparty: entities.Account{Name: "employee"},
				Amount:       decimal.New(1000, 0),
			},
		}
		storage.EXPECT().GetPaymentsList(ctx).Return(storageResult, nil)

		payments, err := banking.NewService(storage).GetPaymentsList(ctx)
		require.NoError(t, err)
		assert.Equal(t, storageResult, payments)
	})

	t.Run("propagates storage exceptions", func(t *testing.T) {
		mCtrl := gomock.NewController(t)
		defer mCtrl.Finish()
		storage := mocks.NewMockStorage(mCtrl)

		storage.EXPECT().GetPaymentsList(ctx).Return(nil, ErrDB)

		_, err := banking.NewService(storage).GetPaymentsList(ctx)
		require.Error(t, err)
	})
}
func TestBankingSvcSendPayment(t *testing.T) {
	t.Run("creates payment silently", func(t *testing.T) {
		mCtrl := gomock.NewController(t)
		defer mCtrl.Finish()
		storage := mocks.NewMockStorage(mCtrl)

		from := entities.Account{Name: "sender"}
		to := entities.Account{Name: "receiver"}
		amount := decimal.New(1356, -2)
		storage.EXPECT().SendPayment(ctx, from, to, amount).Return(nil)

		err := banking.NewService(storage).SendPayment(ctx, from, to, amount)
		require.NoError(t, err)
	})

	t.Run("catches transfer attempts to same account", func(t *testing.T) {
		mCtrl := gomock.NewController(t)
		defer mCtrl.Finish()
		storage := mocks.NewMockStorage(mCtrl)

		sender := entities.Account{Name: "sender"}
		amount := decimal.New(1356, -2)

		err := banking.NewService(storage).SendPayment(ctx, sender, sender, amount)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "can't transfer funds to the same account")
	})

	t.Run("propagates storage exceptions", func(t *testing.T) {
		mCtrl := gomock.NewController(t)
		defer mCtrl.Finish()
		storage := mocks.NewMockStorage(mCtrl)

		from := entities.Account{Name: "sender"}
		to := entities.Account{Name: "receiver"}
		amount := decimal.New(1356, -2)
		storage.EXPECT().SendPayment(ctx, from, to, amount).Return(ErrDB)

		err := banking.NewService(storage).SendPayment(ctx, from, to, amount)
		require.Error(t, err)
	})
}
