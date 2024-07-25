package constant

type OrderType bool

const (
	BuyOrder  OrderType = true
	SellOrder OrderType = false
)

// String returns a string representation of the OrderType.
func (o OrderType) String() string {
	switch o {
	case BuyOrder:
		return "BuyOrder"
	default:
		return "SellOrder"
	}
}
