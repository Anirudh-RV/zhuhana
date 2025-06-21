package runtime

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

var terminateStatus = map[OrderStatus]struct{}{
	models.StatusComplete:  {},
	models.StatusExpired:   {},
	models.StatusError:     {},
	models.StatusCancelled: {},
}

// OrderFlow Map-based FSM
type OrderFlow struct {
	OrderID      string
	OrderRequest *models.OrderRequest
	status       OrderStatus
	history      []OrderStatusTransition
	OnChange     func(from OrderStatus, to OrderStatus) // hook
}

func NewOrderFlow(orderRequest *models.OrderRequest) *OrderFlow {
	return &OrderFlow{
		OrderID:      orderRequest.OrderID,
		OrderRequest: orderRequest,
		status:       models.StatusNew,
		history:      []OrderStatusTransition{},
	}
}

func (fsm *OrderFlow) Current() OrderStatus {
	return fsm.status
}

func (fsm *OrderFlow) GetHistory() []OrderStatusTransition {
	return fsm.history
}

func (fsm *OrderFlow) CanTransition(to OrderStatus) bool {
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

func (fsm *OrderFlow) IsTerminated() bool {
	_, ok := terminateStatus[fsm.status]
	return ok
}

func (fsm *OrderFlow) Transition(to OrderStatus) error {
	if !fsm.CanTransition(to) {
		return fmt.Errorf("orderflow invalid transition: %s → %s", fsm.status, to)
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
