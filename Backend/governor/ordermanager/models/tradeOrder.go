package models

// TradeOrderSide: side of trade order.
type TradeOrderSide string
const (
    SideBuy  TradeOrderSide = "BUY"
    SideSell TradeOrderSide = "SELL"
	SideShort TradeOrderSide = "SHORT" // Represents a short sell order
    SideInvalid TradeOrderSide = "INVALID" // Placeholder for an invalid side
)

type TradeOrderType string
const (
    TypeMarket  TradeOrderType = "MARKET"
    TypeLimit   TradeOrderType = "LIMIT" 
    TypeStop    TradeOrderType = "STOP"
    TypeStopLimit TradeOrderType = "STOP_LIMIT"
    TypeTrailingStop TradeOrderType = "TRAILING_STOP"
    TypeFillOrKill TradeOrderType = "FILL_OR_KILL"
    TypeImmediateOrCancel TradeOrderType = "IMMEDIATE_OR_CANCEL"
    TypeAllOrNone TradeOrderType = "ALL_OR_NONE"
    TypeInvalid TradeOrderType = "INVALID" // Placeholder for an invalid type
)

// TradeOrder represents a trade order in the system.
// It includes fields for the symbol, side, quantity, price, and type of the order.
type TradeOrder struct {
    Symbol      string `json:"symbol"` // Stock symbol, e.g., "AAPL", "GOOGL"
    Side        TradeOrderSide `json:"side"`
    Type        TradeOrderType `json:"type"` // Type of order, e.g., "MARKET", "LIMIT"
    Quantity    float64 `json:"quantity"` // Number of shares
    Price       float64 `json:"price"` // Price per share for limit orders
    Domain      string `json:"domain"`// Trading domain, e.g., "NASDAQ"
}

