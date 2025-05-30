package models

type TradeOrderRequest struct {
	AuthToken  string `json:"authToken"` // Authentication token for the user
    UserID       string `json:"userId"` // Unique identifier for User
    StrategyID string `json:"strategyId"` // Unique identifier for the strategy
    OrderID    string `json:"orderId"` // Unique identifier for the order
    Timestamp  time.Time `json:"timeStamp"` // Timestamp of the order request
    Priority   int  `json:"priority"` // Priority of the order, lower values indicate higher priority
    Order      TradeOrder `json:"order"` // TradeOrder contains details about the trade order
    // Metadata       interface{} `json:"meta,omitempty"`// Metadata related to the order, e.g., user notes, tags, or additional attributes for processing
}