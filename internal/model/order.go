package model

import "time"

type Order struct {
	ID         uint64
	CustomerID uint
	Price      uint
	Timestamp  time.Time
	GTT        *time.Time // Good Til Time
}
