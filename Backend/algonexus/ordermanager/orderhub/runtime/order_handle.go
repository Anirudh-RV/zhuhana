package runtime

import (
	rsm "algonexus/eventqueue/models"
	"algonexus/ordermanager/models"
)

type OrderHandle struct {
	OrderID   string
	OrderFlow *OrderFlow
	Channel   chan rsm.RsEvent
}

func NewOrderHandle(req *models.OrderRequest) *OrderHandle {
	return &OrderHandle{
		OrderID:   req.OrderID,
		OrderFlow: NewOrderFlow(req),
		Channel:   make(chan rsm.RsEvent),
	}
}
