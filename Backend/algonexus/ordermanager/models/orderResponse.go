package models

import "time"

// OrderStatus Strong Type-style "Enum" Check
// TODO: json decode check as well

type OrderResponse struct {
	OrderID       string      `json:"order_id"`
	OrderDetails  Order       `json:"order_details"` // This field is only for debugging
	SubmitTime    time.Time   `json:"submit_time"`
	BrokerOrderID string      `json:"broker_order_id"`
	Status        OrderStatus `json:"status"` // e.g. "filled", "pending", "rejected"
	Message       string      `json:"message"`
	Fills         []OrderFill `json:"fills"`
	RejectReason  string      `json:"reject_reason,omitempty"`
	Time          time.Time   `json:"time"`
}
