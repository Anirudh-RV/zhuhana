package consumers

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

type RsOrderConsumerMsgHandler struct{}

func (h *RsOrderConsumerMsgHandler) Handle(ctx context.Context, msg redis.XMessage) error {

	if msg.ID == "" {
		return fmt.Errorf("invalid msg id")
	}
	data, ok := msg.Values["data"]
	if !ok {
		return fmt.Errorf("missing data in message %s", msg.ID)
	}

	fmt.Printf("Wow, order consumed: %s", data.(string))

	return nil
}
