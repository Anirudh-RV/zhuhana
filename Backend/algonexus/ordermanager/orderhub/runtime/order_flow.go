package runtime

import (
	"algonexus/ordermanager/models"
	"fmt"
	"sync"
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

// OrderFlow is a map-based FSM.
//
// Concurrency: across one order's life the producer (ingress goroutine), a submit
// worker (anchor pool) and a Listener (fills pool) advance the same OrderFlow from
// different goroutines, so all access to status/history is guarded by mu.
type OrderFlow struct {
	OrderID      string
	OrderRequest *models.OrderRequest
	mu           sync.Mutex
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
	fsm.mu.Lock()
	defer fsm.mu.Unlock()
	return fsm.status
}

func (fsm *OrderFlow) GetHistory() []OrderStatusTransition {
	fsm.mu.Lock()
	defer fsm.mu.Unlock()
	out := make([]OrderStatusTransition, len(fsm.history))
	copy(out, fsm.history)
	return out
}

func (fsm *OrderFlow) CanTransition(to OrderStatus) bool {
	fsm.mu.Lock()
	defer fsm.mu.Unlock()
	return canTransition(fsm.status, to)
}

// canTransition is the lock-free core, callable while mu is already held.
func canTransition(from OrderStatus, to OrderStatus) bool {
	allowed, ok := validTransitions[from]
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
	fsm.mu.Lock()
	defer fsm.mu.Unlock()
	_, ok := terminateStatus[fsm.status]
	return ok
}

func (fsm *OrderFlow) Transition(to OrderStatus) error {
	fsm.mu.Lock()
	if !canTransition(fsm.status, to) {
		from := fsm.status
		fsm.mu.Unlock()
		return fmt.Errorf("orderflow invalid transition: %s → %s", from, to)
	}
	from := fsm.status
	fsm.status = to
	fsm.history = append(fsm.history, OrderStatusTransition{
		From:      from,
		To:        to,
		Timestamp: time.Now(),
	})
	hook := fsm.OnChange
	fsm.mu.Unlock()

	if hook != nil {
		hook(from, to) // invoked without the lock held
	}
	return nil
}
