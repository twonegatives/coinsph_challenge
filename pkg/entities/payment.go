package entities

import "github.com/shopspring/decimal"

// Payment represents a money move between two accounts in a single direction
// (either incoming or outgoing). Each payment has a corresponding opposite payment.
// Both such payments are linked by a single Transaction.
type Payment struct {
	Account      Account
	Counterparty Account
	Transaction  Transaction
	Direction    Direction
	Amount       decimal.Decimal
	Currency     Currency
}
