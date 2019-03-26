package entities

import "github.com/shopspring/decimal"

type Payment struct {
	Account     string          `json:"account"`
	Amount      decimal.Decimal `json:"amount"`
	FromAccount string          `json:"from_account"`
	ToAccount   string          `json:"to_account"`
	Direction   Direction       `json:"direction"`
}
