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
	TypeFillOrKill        OrderType = "FILL_OR_KILL"
	TypeImmediateOrCancel OrderType = "IMMEDIATE_OR_CANCEL"
	TypeAllOrNone         OrderType = "ALL_OR_NONE"
	TypeInvalid           OrderType = "INVALID" // Placeholder for an invalid type
)

type OrderMode string

const (
	ModeIntraday OrderMode = "INTRADAY"
	ModeDelivery OrderMode = "DELIVERY"
)

type OrderDomain string

const (
	DomainBacktest OrderDomain = "BACKTEST"
)

type OrderTIF string

const (
	TIFDay OrderTIF = "DAY"
	TIFGtc OrderTIF = "GTC"
	TIFIoc OrderTIF = "IOC"
)

// Order represents a trade order in the system.
// It includes fields for the symbol, side, quantity, price, and type of the order.
type Order struct {
	Symbol      string      `json:"symbol" binding:"required"` // Stock symbol, e.g., "AAPL", "GOOGL"
	Mode        OrderMode   `json:"mode" binding:"required"`
	Side        OrderSide   `json:"side" binding:"required"`
	Type        OrderType   `json:"type" binding:"required"` // Type of order, e.g., "MARKET", "LIMIT"
	Domain      OrderDomain `json:"domain"`                  // Trading domain, e.g., "NASDAQ"
	TimeInForce OrderTIF    `json:"time_in_force"`
	Quantity    float64     `json:"quantity" binding:"required"` // Number of shares
	Price       float64     `json:"price" binding:"required" `   // Price per share for limit orders
	Priority    int         `json:"priority"`                    // Priority of the order, lower values indicate higher priority
}
