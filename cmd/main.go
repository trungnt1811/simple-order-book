package main

import (
	"fmt"

	"github.com/trungnt1811/simple-order-book/internal/module"
)

func main() {
	ob := module.NewOrderBook()

	// Example usage
	ob.SubmitOrder(1, 100, true, nil)  // Customer 1 offers to buy at $100
	ob.SubmitOrder(2, 99, true, nil)   // Customer 2 offers to buy at $99
	ob.SubmitOrder(3, 110, false, nil) // Customer 3 offers to sell at $110
	ob.SubmitOrder(4, 105, false, nil) // Customer 4 offers to sell at $105

	// Example query
	orders := ob.QueryOrders(1)
	fmt.Println("Customer 1's active orders:")
	for _, order := range orders {
		fmt.Printf("Order ID: %d, Price: %d\n", order.ID, order.Price)
	}

	// Example cancellation
	ob.CancelOrder(1)
}
