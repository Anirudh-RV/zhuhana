package runtime

import (
	"algonexus/ordermanager/models"
	hubmodels "algonexus/ordermanager/orderhub/models"
)

type OrderSession struct {
	OrderID   string
	OrderFlow *OrderFlow
	Channel   chan hubmodels.OrderEvent
}

func NewOrderSession(req *models.OrderRequest) *OrderSession {
	return &OrderSession{
		OrderID:   req.OrderID,
		OrderFlow: NewOrderFlow(req),
		Channel:   make(chan hubmodels.OrderEvent),
	}
}

// TODO OrderSession methods

// func (h *OrderSession) Submit() error { return nil }
//
// func (h *OrderSession) Cancel() error { return nil }
//
// func (h *OrderSession) HandleBrokerResponse(res *models.OrderResponse) error { return nil }
//
// func (h *OrderSession) OnChangeStatus(hook func(from, to OrderStatus)) {}

func (ctx *OrderSession) ApplyEvent(e hubmodels.OrderEvent) {
	// delegate to FSM
	// TODO
}
