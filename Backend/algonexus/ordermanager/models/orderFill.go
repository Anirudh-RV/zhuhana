package models

import "time"

type OrderFill struct {
	FillID   string    `json:"fill_id"`  // uuid for each transaction
	OrderID  string    `json:"order_id"` // OrderID it belongs to
	Price    float64   `json:"price"`
	Quantity float64   `json:"quantity"`
	Side     string    `json:"side"`
	Source   string    `json:"source"` // adapter/broker name e.g. "simulator", "alpaca"
	Time     time.Time `json:"time"`
	Status   string    `json:"status"` // broker provided
}
