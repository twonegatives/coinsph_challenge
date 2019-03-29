package entities

import "time"

type Transaction struct {
	ID        int       `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}
