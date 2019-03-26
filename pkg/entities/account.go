package entities

import "github.com/shopspring/decimal"

type Account struct {
	ID       string          `json:"id"`
	Balance  decimal.Decimal `json:"balance"`
	Currency Currency        `json:"currency"`
}
