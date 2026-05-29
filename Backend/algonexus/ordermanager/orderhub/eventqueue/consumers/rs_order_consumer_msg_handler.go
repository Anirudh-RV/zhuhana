package consumers

import (
	"algonexus/logger"
	"algonexus/ordermanager/models"
	"algonexus/ordermanager/orderhub/ports"
	"algonexus/ordermanager/orderhub/registry"
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// RsOrderConsumerMsgHandler is the order-system anchor (the submit worker). It takes an
// order off the ingress stream, advances the FSM ASSIGNED -> SUBMITTED, and hands it to
// the broker via the in-process Broker port. The fill comes back asynchronously on
// broker.Fills() and is finalized by the OrderHub Listener pool.
type RsOrderConsumerMsgHandler struct {
	logger   *logger.Logger
	registry *registry.OrderHubRegistry
	broker   ports.Broker
}

func NewRsOrderConsumerMsgHandler(logger *logger.Logger, registry *registry.OrderHubRegistry, broker ports.Broker) *RsOrderConsumerMsgHandler {
	return &RsOrderConsumerMsgHandler{
		logger:   logger,
		registry: registry,
		broker:   broker,
	}
}

func (h *RsOrderConsumerMsgHandler) Handle(ctx context.Context, msg redis.XMessage) error {
	if msg.ID == "" {
		return fmt.Errorf("invalid msg id")
	}
	data, ok := msg.Values["data"]
	if !ok {
		return fmt.Errorf("missing data in message %s", msg.ID)
	}

	var req models.OrderRequest
	if err := json.Unmarshal([]byte(data.(string)), &req); err != nil {
		return fmt.Errorf("unmarshal order request: %w", err)
	}

	handle := h.registry.Get(req.OrderID)
	if handle == nil {
		h.logger.Warning("anchor: no handle for order, skipping", zap.String("orderID", req.OrderID))
		return nil
	}

	// Worker picked it up (ASSIGNED), then submits to the broker (SUBMITTED).
	// SUBMITTED is set BEFORE handing to the broker so the async fill can never arrive
	// before the FSM reaches SUBMITTED (which would be an illegal CONFIRMED transition).
	h.transition(req.OrderID, handle.OrderFlow.Transition(models.StatusAssigned), models.StatusAssigned)
	h.transition(req.OrderID, handle.OrderFlow.Transition(models.StatusSubmitted), models.StatusSubmitted)

	if err := h.broker.Submit(ctx, req); err != nil {
		// fail-fast: broker did not accept the order -> terminal ERROR + cleanup here.
		// The broker never received it, so no fill will come and the Listener pool will
		// never touch this order — no double cleanup.
		h.transition(req.OrderID, handle.OrderFlow.Transition(models.StatusError), models.StatusError)
		h.registry.Delete(req.OrderID)
		h.logger.Warning("broker rejected submit (fail-fast)", zap.String("orderID", req.OrderID), zap.Error(err))
		return nil
	}

	h.logger.Info("anchor submitted order to broker", zap.String("orderID", req.OrderID))
	return nil
}

func (h *RsOrderConsumerMsgHandler) transition(orderID string, err error, to models.OrderStatus) {
	if err != nil {
		h.logger.Error("FSM transition failed",
			zap.String("orderID", orderID),
			zap.String("to", string(to)),
			zap.Error(err))
	}
}
