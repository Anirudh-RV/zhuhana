package consumers

import (
	"algonexus/logger"
	"algonexus/ordermanager/models"
	hubModels "algonexus/ordermanager/orderhub/models"
	"algonexus/ordermanager/orderhub/registry"
	"encoding/json"
	"go.uber.org/zap"
	"time"

	//orderHubServices "algonexus/ordermanager/orderhub/services"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

type RsOrderConsumerMsgHandler struct {
	logger   *logger.Logger
	registry *registry.OrderHubRegistry
}

func NewRsOrderConsumerMsgHandler(logger *logger.Logger, registry *registry.OrderHubRegistry) *RsOrderConsumerMsgHandler {
	return &RsOrderConsumerMsgHandler{
		registry: registry,
		logger:   logger,
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
		panic(err)
	}

	h.logger.Info("Wow, order consumed", zap.String("data", data.(string)))

	session := h.registry.Get(req.OrderID)
	if session == nil {
		h.logger.Error("couldn't find the handle", zap.String("orderID", req.OrderID))
		return fmt.Errorf("couldn't find the handle %s", req.OrderID)
	}

	h.logger.Info("OrderHub session Lookup check", zap.String("order status", string(session.OrderFlow.Current())))

	//TODO backtest
	//Test event
	go func() {
		event := &hubModels.OrderEvent{
			OrderID:   req.OrderID,
			Timestamp: time.Now().Unix(),
			Type:      hubModels.EventBrokerConfirmed,
			Payload:   nil,
		}

		session.Channel <- *event
	}()

	return nil
}
