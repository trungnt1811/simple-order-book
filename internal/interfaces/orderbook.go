package interfaces

import (
	"time"

	"github.com/trungnt1811/simple-order-book/internal/constant"
	"github.com/trungnt1811/simple-order-book/internal/model"
)

type OrderBookUCase interface {
	SubmitOrder(customerID uint, price uint, orderType constant.OrderType, gtt *time.Time) error
	CancelOrder(orderID uint64) error
	QueryOrders(customerID uint) []*model.Order
	GetNextOrderID() uint64
	GetSellOrders() model.OrderHeap
	GetBuyOrders() model.OrderHeap
	GetOrders() map[uint64]*model.Order
	GetCustomerOrders() map[uint]map[uint64]*model.Order
}
