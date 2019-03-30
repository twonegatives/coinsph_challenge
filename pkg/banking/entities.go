package banking

import (
	"github.com/shopspring/decimal"
	"github.com/twonegatives/coinsph_challenge/pkg/entities"
)

// sendPaymentRequest is a structure which banking transport layer
// uses to pass data further to endpoint layer on
// POST /api/v1/payments request
type sendPaymentRequest struct {
	From   entities.Account
	To     entities.Account
	Amount decimal.Decimal
}

// createAccountRequest is a structure which banking transport layer
// uses to pass data further to endpoint layer on
// POST /api/v1/accounts request
type createAccountRequest struct {
	Name string
}

// getAccountsResponse is a structure which banking endpoint layer
// uses to pass data upside down to transport layer on
// GET /api/v1/accounts
type getAccountsResponse struct {
	Accounts []entities.Account `json:"accounts"`
}

// getPaymentsResponse is a structure which banking endpoint layer
// uses to pass data upside down to transport layer on
// GET /api/v1/payments
type getPaymentsResponse struct {
	Payments []entities.Payment `json:"payments"`
}

// createAccountsResponse is a structure which banking endpoint layer
// uses to pass data upside down to transport layer on
// POST /api/v1/accounts
type createAccountResponse struct {
	Account entities.Account `json:"account"`
}
