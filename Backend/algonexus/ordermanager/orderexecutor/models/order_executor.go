package models

import "algonexus/ordermanager/models"

type OrderExecutor interface {
	Submit(order *models.OrderRequest) (*OrderAck, error)
	Subscribe(orderID string) (<-chan *OrderStatus, error)
}
