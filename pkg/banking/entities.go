package banking

import (
	"github.com/shopspring/decimal"
	"github.com/twonegatives/coinsph_challenge/pkg/entities"
)

type sendPaymentRequest struct {
	From   entities.Account
	To     entities.Account
	Amount decimal.Decimal
}

type createAccountRequest struct {
	Name string
}

type getAccountsResponse struct {
	Accounts []entities.Account `json:"accounts"`
}

type getPaymentsResponse struct {
	Payments []entities.Payment `json:"payments"`
}

type createAccountResponse struct {
	Account entities.Account `json:"account"`
}
