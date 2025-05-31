package models

// OrderSide side of trade order.
type OrderSide string

const (
	SideBuy     OrderSide = "BUY"
	SideSell    OrderSide = "SELL"
	SideShort   OrderSide = "SHORT"   // Represents a short sell order
	SideInvalid OrderSide = "INVALID" // Placeholder for an invalid side
)

type OrderType string

const (
	TypeMarket            OrderType = "MARKET"
	TypeLimit             OrderType = "LIMIT"
	TypeStop              OrderType = "STOP"
	TypeStopLimit         OrderType = "STOP_LIMIT"
	TypeTrailingStop      OrderType = "TRAILING_STOP"
	TypeFillOrKill        OrderType = "FILL_OR_KILL"
	TypeImmediateOrCancel OrderType = "IMMEDIATE_OR_CANCEL"
	TypeAllOrNone         OrderType = "ALL_OR_NONE"
	TypeInvalid           OrderType = "INVALID" // Placeholder for an invalid type
)

// Order represents a trade order in the system.
// It includes fields for the symbol, side, quantity, price, and type of the order.
type Order struct {
	Symbol   string    `json:"symbol"` // Stock symbol, e.g., "AAPL", "GOOGL"
	Side     OrderSide `json:"side"`
	Type     OrderType `json:"type"`     // Type of order, e.g., "MARKET", "LIMIT"
	Quantity float64   `json:"quantity"` // Number of shares
	Price    float64   `json:"price"`    // Price per share for limit orders
	Domain   string    `json:"domain"`   // Trading domain, e.g., "NASDAQ"
}
