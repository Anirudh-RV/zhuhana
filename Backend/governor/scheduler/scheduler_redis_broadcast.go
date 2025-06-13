package scheduler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (scs *SchedulerService) BroadcastCancelJob(UserAlgorithmID uuid.UUID, cronEntryID int64) error {
	ctx := context.Background()
	payload := map[string]any{
		"user_algorithm_id": UserAlgorithmID.String(),
		"cron_entry_id":     cronEntryID,
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	scs.logger.Info("publishing cron_cancel_channel on redis", zap.String("execution level", "BroadcastCancelJob"), zap.String("UserAlgorithmID", UserAlgorithmID.String()), zap.String("cronEntryID", fmt.Sprint(cronEntryID)))

	return scs.redisObj.Publish(ctx, "cron_cancel_channel", data).Err()
}
