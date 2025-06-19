package models

type OrderStatus string

const (
	StatusNew             OrderStatus = "NEW"
	StatusPendingSend     OrderStatus = "PENDING_SEND" // Ready to be enqueued
	StatusEnqueued        OrderStatus = "ENQUEUED"     // In the event queue
	StatusAssigned        OrderStatus = "ASSIGNED"     // Dequeued to a worker
	StatusSubmitted       OrderStatus = "SUBMITTED"    //Submitted to broker
	StatusRejected        OrderStatus = "REJECTED"
	StatusBrokerConfirmed OrderStatus = "CONFIRMED"      // Got broker response
	StatusInTransaction   OrderStatus = "IN_TRANSACTION" // During the transaction
	StatusComplete        OrderStatus = "COMPLETED"      // Order fully completed

	StatusCancelled OrderStatus = "CANCELLED"
	StatusExpired   OrderStatus = "EXPIRED"
	StatusError     OrderStatus = "ERROR" // Internal Error
	StatusUnknown   OrderStatus = "UNKNOWN"

	StatusBrokerRejected        OrderStatus = "BROKER_REJECTED"
	StatusBrokerPartiallyFilled OrderStatus = "BROKER_PARTIALLY_FILLED"
	StatusBrokerFilled          OrderStatus = "BROKER_FILLED"
)
