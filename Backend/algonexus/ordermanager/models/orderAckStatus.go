package models

type OrderAckStatus string

const (
	AckStatusSubmitted OrderAckStatus = "SUBMITTED"
	ActStatusRejected  OrderAckStatus = "REJECT"
)
