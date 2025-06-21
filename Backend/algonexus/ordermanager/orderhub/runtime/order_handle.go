package runtime

import (
	"algonexus/ordermanager/models"
	hubmodels "algonexus/ordermanager/orderhub/models"
)

type OrderHandle struct {
	OrderID   string
	OrderFlow *OrderFlow
	Channel   chan hubmodels.OrderEvent
}

func NewOrderHandle(req *models.OrderRequest) *OrderHandle {
	return &OrderHandle{
		OrderID:   req.OrderID,
		OrderFlow: NewOrderFlow(req),
		Channel:   make(chan hubmodels.OrderEvent),
	}
}
