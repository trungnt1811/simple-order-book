package main

import (
	"sync"

	"go.uber.org/zap"

	"github.com/trungnt1811/simple-order-book/internal/constant"
	"github.com/trungnt1811/simple-order-book/internal/module"
	"github.com/trungnt1811/simple-order-book/internal/util"
)

func main() {
	logger := util.SetupLogger()
	defer logger.Sync() // Flushes buffer, if any

	orderBook := module.NewOrderBookUCase(logger)

	var wg sync.WaitGroup

	// Function to submit multiple orders concurrently
	submitOrders := func(customerID uint, prices []uint, orderType constant.OrderType) {
		defer wg.Done()
		for _, price := range prices {
			orderBook.SubmitOrder(customerID, price, orderType, util.CreateGTT(1))
			logger.Debug("Order submitted", zap.Uint("CustomerID", customerID), zap.Uint("Price", price), zap.String("OrderType", orderType.String()))
		}
	}

	// Start concurrent submissions
	wg.Add(3)
	go submitOrders(1, []uint{110, 109, 108}, constant.BuyOrder)
	go submitOrders(2, []uint{99, 98, 97}, constant.BuyOrder)
	go submitOrders(3, []uint{100, 101, 102}, constant.SellOrder)

	// Wait for all submissions to complete
	wg.Wait()

	// Function to query and log orders
	queryOrdersByCustomerID := func(customerID uint) {
		orders := orderBook.QueryOrders(customerID)
		logger.Debug("Queried active orders", zap.Uint("CustomerID", customerID), zap.Int("OrderCount", len(orders)))
		for _, order := range orders {
			logger.Info("Order details", zap.Uint("CustomerID", customerID), zap.Uint64("OrderID", order.ID), zap.Uint("Price", order.Price), zap.String("OrderType", order.OrderType.String()))
		}
	}

	queryOrdersByCustomerID(1)
	queryOrdersByCustomerID(2)
	queryOrdersByCustomerID(3)

	// Function to cancel multiple orders concurrently
	cancelOrders := func(orderIDs []uint64) {
		defer wg.Done()
		for _, orderID := range orderIDs {
			err := orderBook.CancelOrder(orderID)
			if err != nil {
				logger.Error("Failed to cancel order", zap.Uint64("OrderID", orderID), zap.Error(err))
			} else {
				logger.Debug("Order canceled", zap.Uint64("OrderID", orderID))
			}
		}
	}

	// Cancel some orders concurrently
	wg.Add(2)
	go cancelOrders([]uint64{1, 2, 3, 4})
	go cancelOrders([]uint64{5, 6, 7})

	// Wait for all cancellations to complete
	wg.Wait()

	// Query orders again after cancellation
	queryOrdersByCustomerID(1)
	queryOrdersByCustomerID(2)
	queryOrdersByCustomerID(3)
}
