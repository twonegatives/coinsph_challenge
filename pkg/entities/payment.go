package entities

import "github.com/shopspring/decimal"

type Payment struct {
	Account      Account
	Counterparty Account
	Transaction  Transaction
	Direction    Direction
	Amount       decimal.Decimal
}
