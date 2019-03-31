// Package entities contains global level objects accessible by other
// packages of the application.
package entities

import "github.com/shopspring/decimal"

// Account represents a user account in the system.
type Account struct {
	ID       int             `json:"-"`
	Name     string          `json:"name"`
	Balance  decimal.Decimal `json:"balance"`
	Currency Currency        `json:"currency"`
}

func (a Account) MayGoBelowZero() bool {
	if a.Name == "SYSTEM" {
		return true
	}
	return false
}
