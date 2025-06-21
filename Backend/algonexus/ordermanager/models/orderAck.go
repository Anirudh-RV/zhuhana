package models

import "time"

type OrderAck struct {
	OrderID       string
	ClientOrderID string
	Order         Order
	Time          time.Time //time when received
	Status        OrderAckStatus
}
