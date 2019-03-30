package entities

import "time"

// Transaction is an object linking two related and opposite payments.
type Transaction struct {
	ID        int       `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}
