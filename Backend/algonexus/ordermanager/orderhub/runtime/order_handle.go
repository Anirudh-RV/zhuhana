package runtime

import (
	"algonexus/ordermanager/models"
)

// OrderHandle is the per-order entry in the in-memory pool (the registry). It holds
// the order's FSM. The per-order event channel was removed: fills now come back on the
// broker's single shared fill channel and are finalized by the OrderHub Listener pool,
// so there is no longer a goroutine + channel per in-flight order.
type OrderHandle struct {
	OrderID   string
	OrderFlow *OrderFlow
}

func NewOrderHandle(req *models.OrderRequest) *OrderHandle {
	return &OrderHandle{
		OrderID:   req.OrderID,
		OrderFlow: NewOrderFlow(req),
	}
}
