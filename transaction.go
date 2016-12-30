package gomerchant

import "time"

type Transaction struct {
	ID        string
	Amount    uint
	Currency  string
	Captured  bool
	Paid      bool // if authorized or captured
	Cancelled bool
	Status    string
	CreatedAt *time.Time
	Params
}
