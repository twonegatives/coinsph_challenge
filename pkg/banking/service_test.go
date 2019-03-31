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

type account struct {
	name          string
	balanceBefore decimal.Decimal
	balanceAfter  decimal.Decimal
}

type paymentUsecase struct {
	sender   account
	receiver account
	amount   decimal.Decimal
	title    string
}

func TestBankingSvcSendPayment(t *testing.T) {
	validTransferCases := []paymentUsecase{
		{
			sender:   account{name: "sender", balanceBefore: decimal.New(1500, 0), balanceAfter: decimal.New(148450, -2)},
			receiver: account{name: "receiver", balanceBefore: decimal.New(0, 0), balanceAfter: decimal.New(1550, -2)},
			amount:   decimal.New(1550, -2),
			title:    "for transfer between user accounts",
		},
		{
			sender:   account{name: "SYSTEM", balanceBefore: decimal.New(0, 0), balanceAfter: decimal.New(-150, 0)},
			receiver: account{name: "receiver", balanceBefore: decimal.New(0, 0), balanceAfter: decimal.New(150, 0)},
			amount:   decimal.New(150, 0),
			title:    "for transfer of SYSTEM with overdraft (allows going below zero)",
		},
	}

	t.Run("creates payment silently", func(t *testing.T) {
		for _, tt := range validTransferCases {
			t.Run(tt.title, func(t *testing.T) {
				mCtrl := gomock.NewController(t)
				defer mCtrl.Finish()
				storage := mocks.NewMockStorage(mCtrl)

				from := entities.Account{Name: tt.sender.name, Balance: tt.sender.balanceBefore}
				to := entities.Account{Name: tt.receiver.name, Balance: tt.receiver.balanceBefore}
				amount := tt.amount

				newSender := entities.Account{
					Name:    from.Name,
					Balance: tt.sender.balanceAfter,
				}

				newReceiver := entities.Account{
					Name:    to.Name,
					Balance: tt.receiver.balanceAfter,
				}

				outgoing := entities.Payment{
					Account:      from,
					Counterparty: to,
					Amount:       amount,
					Transaction:  entities.Transaction{},
					Direction:    entities.Outgoing,
					Currency:     entities.USD,
				}

				incoming := entities.Payment{
					Account:      to,
					Counterparty: from,
					Amount:       amount,
					Transaction:  entities.Transaction{},
					Direction:    entities.Incoming,
					Currency:     entities.USD,
				}

				storage.EXPECT().BeginTx(gomock.Any(), nil).Return(storage, nil)
				storage.EXPECT().GetAccountForUpdate(gomock.Any(), &from).Return(nil)
				storage.EXPECT().GetAccountForUpdate(gomock.Any(), &to).Return(nil)
				storage.EXPECT().CreateTransaction(gomock.Any()).Return(entities.Transaction{}, nil)
				storage.EXPECT().SendPayment(gomock.Any(), outgoing).Return(nil)
				storage.EXPECT().SendPayment(gomock.Any(), incoming).Return(nil)
				storage.EXPECT().SetAccountBalance(gomock.Any(), newSender).Return(nil)
				storage.EXPECT().SetAccountBalance(gomock.Any(), newReceiver).Return(nil)
				storage.EXPECT().CommitTx(gomock.Any()).Return(nil)
				storage.EXPECT().RollbackTx(gomock.Any()).Return(nil)

				err := banking.NewService(storage).SendPayment(ctx, from, to, amount)
				require.NoError(t, err)
			})
		}

	})

	t.Run("catches transfer attempts to same account", func(t *testing.T) {
		mCtrl := gomock.NewController(t)
		defer mCtrl.Finish()
		storage := mocks.NewMockStorage(mCtrl)

		sender := entities.Account{Name: "benjamin"}
		receiver := entities.Account{Name: "benjamin"}
		amount := decimal.New(1356, -2)

		err := banking.NewService(storage).SendPayment(ctx, sender, receiver, amount)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "can't transfer funds to the same account")
	})

	t.Run("catches transfer attempts exceeding balance", func(t *testing.T) {
		mCtrl := gomock.NewController(t)
		defer mCtrl.Finish()
		storage := mocks.NewMockStorage(mCtrl)

		from := entities.Account{Name: "sender", Balance: decimal.New(15, 0)}
		to := entities.Account{Name: "receiver"}
		amount := decimal.New(50, 0)

		storage.EXPECT().BeginTx(gomock.Any(), nil).Return(storage, nil)
		storage.EXPECT().GetAccountForUpdate(gomock.Any(), &from).Return(nil)
		storage.EXPECT().GetAccountForUpdate(gomock.Any(), &to).Return(nil)
		storage.EXPECT().RollbackTx(gomock.Any()).Return(nil)

		err := banking.NewService(storage).SendPayment(ctx, from, to, amount)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "sender account has insufficient funds")
	})

	t.Run("propagates storage exceptions", func(t *testing.T) {
		from := entities.Account{Name: "sender", Balance: decimal.New(15, 0)}
		to := entities.Account{Name: "receiver"}
		amount := decimal.New(10, 0)

		t.Run("on db tx opening", func(t *testing.T) {
			mCtrl := gomock.NewController(t)
			defer mCtrl.Finish()
			storage := mocks.NewMockStorage(mCtrl)

			storage.EXPECT().BeginTx(gomock.Any(), nil).Return(storage, ErrDB)
			storage.EXPECT().RollbackTx(gomock.Any()).Return(nil)

			err := banking.NewService(storage).SendPayment(ctx, from, to, amount)
			require.Error(t, err)
			assert.Contains(t, err.Error(), "can't open transaction")
		})

		t.Run("on accounts obtaining", func(t *testing.T) {
			mCtrl := gomock.NewController(t)
			defer mCtrl.Finish()
			storage := mocks.NewMockStorage(mCtrl)

			storage.EXPECT().BeginTx(gomock.Any(), nil).Return(storage, nil)
			storage.EXPECT().GetAccountForUpdate(gomock.Any(), &from).Return(ErrDB)
			storage.EXPECT().RollbackTx(gomock.Any()).Return(nil)

			err := banking.NewService(storage).SendPayment(ctx, from, to, amount)
			require.Error(t, err)
			assert.Contains(t, err.Error(), "can't obtain sender account")
		})

		t.Run("on inserting outgoing payment", func(t *testing.T) {
			mCtrl := gomock.NewController(t)
			defer mCtrl.Finish()
			storage := mocks.NewMockStorage(mCtrl)

			storage.EXPECT().BeginTx(gomock.Any(), nil).Return(storage, nil)
			storage.EXPECT().GetAccountForUpdate(gomock.Any(), &from).Return(nil)
			storage.EXPECT().GetAccountForUpdate(gomock.Any(), &to).Return(nil)
			storage.EXPECT().CreateTransaction(gomock.Any()).Return(entities.Transaction{}, nil)
			storage.EXPECT().SendPayment(gomock.Any(), gomock.Any()).Return(ErrDB)
			storage.EXPECT().RollbackTx(gomock.Any()).Return(nil)

			err := banking.NewService(storage).SendPayment(ctx, from, to, amount)
			require.Error(t, err)
			assert.Contains(t, err.Error(), "can't insert outgoing payment")
		})
	})
}
