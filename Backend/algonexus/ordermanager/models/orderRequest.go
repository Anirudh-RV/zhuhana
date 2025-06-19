package models

import "time"

type OrderRequest struct {
	//UserID     string    `json:"userId"`     // Unique identifier for User
	//StrategyID string    `json:"strategyId"` // Unique identifier for the strategy
	OrderID   string    `json:"orderId"`   // Unique identifier for the order
	Order     Order     `json:"order"`     // Order contains details about the trade order
	Timestamp time.Time `json:"timeStamp"` // Timestamp of the order request (client side)
	// Metadata       interface{} `json:"meta,omitempty"`// Metadata related to the order, e.g., user notes, tags, or additional attributes for processing
}
