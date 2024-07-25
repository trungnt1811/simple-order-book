package main

import (
	"fmt"
	"sync"

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
	submitOrders := func(customerID uint, prices []uint, orderType constant.OrderType, wg *sync.WaitGroup) {
		defer wg.Done()
		for _, price := range prices {
			orderBook.SubmitOrder(customerID, price, orderType, util.CreateGTT(1))
		}
	}

	// Function to cancel multiple orders concurrently
	cancelOrders := func(orderIDs []uint64, wg *sync.WaitGroup) {
		defer wg.Done()
		for _, orderID := range orderIDs {
			err := orderBook.CancelOrder(orderID)
			if err != nil {
				fmt.Printf("Failed to cancel order %d: %v\n", orderID, err)
			}
		}
	}

	// Start concurrent submissions
	wg.Add(3)
	go submitOrders(1, []uint{100, 101, 102}, constant.BuyOrder, &wg)
	go submitOrders(2, []uint{99, 98, 97}, constant.BuyOrder, &wg)
	go submitOrders(3, []uint{110, 109, 108}, constant.SellOrder, &wg)

	// Wait for all submissions to complete
	wg.Wait()

	// Query orders
	queryOrders := func(customerID uint) {
		orders := orderBook.QueryOrders(customerID)
		fmt.Printf("Customer %d's active orders:\n", customerID)
		for _, order := range orders {
			fmt.Printf("Order ID: %d, Price: %d, Type: %v\n", order.ID, order.Price, order.OrderType)
		}
	}

	queryOrders(1)
	queryOrders(2)
	queryOrders(3)

	// Cancel some orders concurrently
	orderIDs := []uint64{1, 2, 3, 4, 5, 6}
	wg.Add(1)
	go cancelOrders(orderIDs, &wg)

	// Wait for all cancellations to complete
	wg.Wait()

	// Query orders again after cancellation
	queryOrders(1)
	queryOrders(2)
	queryOrders(3)
}
