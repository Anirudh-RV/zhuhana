package runtime

import (
	"algonexus/ordermanager/models"
	hubmodels "algonexus/ordermanager/orderhub/models"
	"fmt"
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

// func (s *OrderSession) Submit() error { return nil }
//
// func (s *OrderSession) Cancel() error { return nil }
//
// func (s *OrderSession) HandleBrokerResponse(res *models.OrderResponse) error { return nil }
//
// func (s *OrderSession) OnChangeStatus(hook func(from, to OrderStatus)) {}

func (s *OrderSession) ApplyEvent(e *hubmodels.OrderEvent) error {
	// delegate to FSM

	if e == nil {
		return fmt.Errorf("unknown event")
	}
	var next models.OrderStatus
	var err error

	//TODO More events
	switch e.Type {
	case hubmodels.EventBrokerConfirmed:
		next = models.StatusBrokerConfirmed
	case hubmodels.EventError:
		next = models.StatusError
	default:
		err = fmt.Errorf("unknown event type %s", e.Type)
	}

	err = s.OrderFlow.Transition(next)
	return err

}
