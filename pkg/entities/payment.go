package entities

import "github.com/shopspring/decimal"

type Payment struct {
	Account     Account
	Participant Account
	Transaction Transaction
	Direction   Direction
	Amount      decimal.Decimal
}
