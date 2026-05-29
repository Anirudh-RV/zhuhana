package ports

import (
	"context"

	brokermodels "algonexus/ordermanager/backtestengine/broker/models"
	"algonexus/ordermanager/models"
)

// Broker is the OrderHub's port to execution (hexagonal). OrderHub depends only on
// this interface; the concrete in-process implementation lives in backtestengine/broker
// and is injected at wiring time (main.go), so backtestengine/broker does NOT import
// orderhub — no import cycle.
//
//   - Submit is a COMMAND: hand an order to the broker and get a synchronous ack/err
//     (fail-fast). It does not wait for the fill.
//   - Fills is an EVENT stream: the broker pushes execution results asynchronously and
//     the OrderHub Listener pool consumes them to finalize the FSM.
type Broker interface {
	Submit(ctx context.Context, req models.OrderRequest) error
	Fills() <-chan brokermodels.BrokerFillEvent
}
