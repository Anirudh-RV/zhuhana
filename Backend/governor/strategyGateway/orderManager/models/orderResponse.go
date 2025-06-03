package models

import "time"

// OrderResponseStatus Strong Type-style "Enum" Check
// TODO: json decode check as well
type OrderResponseStatus string

const (
	ResponseStatusNew             OrderResponseStatus = "NEW"
	ResponseStatusAccepted        OrderResponseStatus = "ACCEPTED"
	ResponseStatusSubmitted       OrderResponseStatus = "SUBMITTED"
	ResponseStatusPartiallyFilled OrderResponseStatus = "PARTIALLY_FILLED"
	ResponseStatusFilled          OrderResponseStatus = "FILLED"
	ResponseStatusRejected        OrderResponseStatus = "FAILED"
	ResponseStatusBrokerRejected  OrderResponseStatus = "BROKER_REJECTED"
	ResponseStatusCancelled       OrderResponseStatus = "CANCELLED"
	ResponseStatusExpired         OrderResponseStatus = "EXPIRED"
	ResponseStatusError           OrderResponseStatus = "ERROR" // Internal Error
)

type OrderResponse struct {
	OrderID       string              `json:"order_id"`
	BrokerOrderID string              `json:"broker_order_id"`
	Status        OrderResponseStatus `json:"status"` // e.g. "filled", "pending", "rejected"
	Message       string              `json:"message"`
	Fills         []OrderFill         `json:"fills"`
	RejectReason  string              `json:"reject_reason,omitempty"`
	Time          time.Time           `json:"time"`
}

type OrderFill struct {
	FillID   string    `json:"fill_id"`  // uuid for each transaction
	OrderID  string    `json:"order_id"` // OrderID it belongs to
	Price    float64   `json:"price"`
	Quantity float64   `json:"quantity"`
	Side     string    `json:"side"`
	Source   string    `json:"source"` // adapter/broker name e.g. "simulator", "alpaca"
	Time     time.Time `json:"time"`
	Status   string    `json:"status"`
}
