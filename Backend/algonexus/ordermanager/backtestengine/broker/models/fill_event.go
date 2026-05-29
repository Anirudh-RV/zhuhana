package models

import ordermodels "algonexus/ordermanager/models"

// BrokerFillEvent is the FillStream payload: a broker's execution result for one
// order, produced by the backtest-engine broker and consumed by the orderhub
// fills consumer (which finalizes the FSM and notifies the listener).
type BrokerFillEvent struct {
	OrderID string                  `json:"order_id"`
	Fill    ordermodels.OrderFill   `json:"fill"`
	Status  ordermodels.OrderStatus `json:"status"`
}
