package consumers

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

type RsOrderConsumerMsgHandler struct{}

func (h *RsOrderConsumerMsgHandler) Handle(ctx context.Context, msg redis.XMessage) error {
	orderIDVal, ok := msg.Values["order_id"]
	if !ok {
		return fmt.Errorf("missing order_id in message %s", msg.ID)
	}
	orderID, ok := orderIDVal.(string)
	if !ok {
		return fmt.Errorf("order_id is not string in message %s", msg.ID)
	}

	return nil
}
