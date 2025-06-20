package runtime

import (
	"algonexus/ordermanager/models"
)

type OrderHandle struct {
	OrderID   string
	OrderFlow *OrderFlow
	Channel   chan models.OrderResponse
}

func NewOrderHandle(req *models.OrderRequest) *OrderHandle {
	return &OrderHandle{
		OrderID:   req.OrderID,
		OrderFlow: NewOrderFlow(req),
		Channel:   make(chan models.OrderResponse),
	}
}
