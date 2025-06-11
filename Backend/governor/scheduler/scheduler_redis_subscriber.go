package scheduler

import (
	"context"
	"encoding/json"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

func (scs *SchedulerService) StartRedisSubscriber(ctx context.Context) {
	pubsub := scs.redisObj.Subscribe(ctx, "cron_cancel_channel")
	scs.logger.Info("started redis subscriber", zap.String("execution level", "StartRedisSubscriber"))

	go func() {
		ch := pubsub.Channel()
		for msg := range ch {
			// Unmarshal the message payload
			var payload struct {
				UserAlgorithmID string `json:"user_algorithm_id"`
				CronEntryID     int64  `json:"cron_entry_id"`
			}
			if err := json.Unmarshal([]byte(msg.Payload), &payload); err != nil {
				scs.logger.Warning("failed to unmarshal cancel payload", zap.Error(err))
				continue
			}

			// Remove the job if it exists
			scs.cronScheduler.Remove(cron.EntryID(payload.CronEntryID))
			scs.logger.Info("🔁 Redis broadcast cancelled job",
				zap.String("JobID", payload.UserAlgorithmID),
				zap.Int64("EntryID", payload.CronEntryID),
			)
		}
	}()
}
