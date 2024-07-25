package model

import (
	"time"

	"github.com/trungnt1811/simple-order-book/internal/constant"
)

type Order struct {
	ID         uint64
	CustomerID uint
	Price      uint
	Timestamp  time.Time
	OrderType  constant.OrderType
	GTT        *time.Time // Good Til Time
}
