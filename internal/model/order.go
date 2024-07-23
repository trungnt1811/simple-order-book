package model

import "time"

type Order struct {
	ID         int
	CustomerID int
	Price      int
	Timestamp  time.Time
	GTT        *time.Time // Good Til Time
}
