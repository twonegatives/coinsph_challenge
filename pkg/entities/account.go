package entities

import "github.com/shopspring/decimal"

type Account struct {
	ID       int             `json:"-"`
	Name     string          `json:"name"`
	Balance  decimal.Decimal `json:"balance"`
	Currency Currency        `json:"currency"`
}
