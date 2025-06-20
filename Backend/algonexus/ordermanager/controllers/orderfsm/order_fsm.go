package orderfsm

import (
	"algonexus/ordermanager/models"
	"fmt"
	"time"
)

type OrderStatus = models.OrderStatus

type OrderStatusTransition struct {
	From      OrderStatus
	To        OrderStatus
	Timestamp time.Time
}

var validTransitions = map[OrderStatus][]OrderStatus{
	models.StatusNew:             {models.StatusPendingSend, models.StatusCancelled, models.StatusExpired, models.StatusError},
	models.StatusPendingSend:     {models.StatusEnqueued, models.StatusCancelled, models.StatusExpired, models.StatusError},
	models.StatusEnqueued:        {models.StatusAssigned, models.StatusCancelled, models.StatusExpired, models.StatusError},
	models.StatusAssigned:        {models.StatusSubmitted, models.StatusRejected, models.StatusCancelled, models.StatusExpired, models.StatusError},
	models.StatusSubmitted:       {models.StatusBrokerConfirmed, models.StatusError},
	models.StatusBrokerConfirmed: {models.StatusInTransaction, models.StatusComplete, models.StatusError},
	models.StatusInTransaction:   {models.StatusBrokerConfirmed, models.StatusError},
}

// OrderFSM Map-based FSM
type OrderFSM struct {
	OrderID      string
	OrderRequest *models.OrderRequest
	status       OrderStatus
	history      []OrderStatusTransition
	OnChange     func(from OrderStatus, to OrderStatus) // hook
}

func NewOrderFSM(orderRequest *models.OrderRequest) *OrderFSM {
	return &OrderFSM{
		OrderID:      orderRequest.OrderID,
		OrderRequest: orderRequest,
		status:       models.StatusNew,
		history:      []OrderStatusTransition{},
	}
}

func (fsm *OrderFSM) Current() OrderStatus {
	return fsm.status
}

func (fsm *OrderFSM) GetHistory() []OrderStatusTransition {
	return fsm.history
}

func (fsm *OrderFSM) CanTransition(to OrderStatus) bool {
	allowed, ok := validTransitions[fsm.status]
	if !ok {
		return false
	}
	for _, s := range allowed {
		if s == to {
			return true
		}
	}
	return false
}

func (fsm *OrderFSM) Transition(to OrderStatus) error {
	if !fsm.CanTransition(to) {
		return fmt.Errorf("orderfsm invalid transition: %s → %s", fsm.status, to)
	}
	from := fsm.status
	fsm.status = to
	fsm.history = append(fsm.history, OrderStatusTransition{
		From:      from,
		To:        to,
		Timestamp: time.Now(),
	})
	if fsm.OnChange != nil {
		fsm.OnChange(from, to) // Hook for all transitions
	}
	return nil
}
