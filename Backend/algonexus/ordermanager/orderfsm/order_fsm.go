package orderfsm

import (
	om "algonexus/ordermanager/models"
	"fmt"
	"time"
)

type OrderStatus = om.OrderStatus

type OrderStatusTransition struct {
	From      OrderStatus
	To        OrderStatus
	Timestamp time.Time
}

var validTransitions = map[OrderStatus][]OrderStatus{
	om.StatusNew:             {om.StatusPendingSend, om.StatusCancelled, om.StatusExpired, om.StatusError},
	om.StatusPendingSend:     {om.StatusEnqueued, om.StatusCancelled, om.StatusExpired, om.StatusError},
	om.StatusEnqueued:        {om.StatusAssigned, om.StatusCancelled, om.StatusExpired, om.StatusError},
	om.StatusAssigned:        {om.StatusSubmitted, om.StatusRejected, om.StatusCancelled, om.StatusExpired, om.StatusError},
	om.StatusSubmitted:       {om.StatusBrokerConfirmed, om.StatusError},
	om.StatusBrokerConfirmed: {om.StatusInTransaction, om.StatusComplete, om.StatusError},
	om.StatusInTransaction:   {om.StatusBrokerConfirmed, om.StatusError},
}

// OrderFSM Map-based FSM
type OrderFSM struct {
	OrderID  string
	status   OrderStatus
	history  []OrderStatusTransition
	OnChange func(from OrderStatus, to OrderStatus) // hook
}

func NewOrderFSM(orderID string) *OrderFSM {
	return &OrderFSM{
		OrderID: orderID,
		status:  om.StatusNew,
		history: []OrderStatusTransition{},
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
		return fmt.Errorf("fsm invalid transition: %s → %s", fsm.status, to)
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
