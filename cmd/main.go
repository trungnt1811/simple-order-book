package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/trungnt1811/simple-order-book/internal/constant"
	"github.com/trungnt1811/simple-order-book/internal/util"
)

func main() {
	ob := util.NewOrderBookWithLogger()

	var wg sync.WaitGroup

	// Function to submit multiple orders concurrently
	submitOrders := func(customerID uint, prices []uint, orderType constant.OrderType, wg *sync.WaitGroup) {
		defer wg.Done()
		for _, price := range prices {
			ob.SubmitOrder(customerID, price, orderType, createGTT(1))
		}
	}

	// Function to cancel multiple orders concurrently
	cancelOrders := func(orderIDs []uint64, wg *sync.WaitGroup) {
		defer wg.Done()
		for _, orderID := range orderIDs {
			err := ob.CancelOrder(orderID)
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
		orders := ob.QueryOrders(customerID)
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

// Helper function to create a GTT time.
func createGTT(hours int) *time.Time {
	gtt := time.Now().Add(time.Duration(hours) * time.Hour)
	return &gtt
}
