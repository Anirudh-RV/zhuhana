package scheduler

import (
	"context"
	"fmt"
	"time"
)

func (scs *SchedulerService) KafkaJobWrapper(job CronJob) func() {
	return func() {
		ctx := context.Background()
		lockKey := fmt.Sprintf("cron-lock:%s", job.ID)
		lock, err := scs.TryLock(ctx, lockKey, time.Minute)
		if err != nil {
			fmt.Println("another node is handling this job:", job.UserAlgorithmID)
			return
		}
		defer lock.Release(ctx)

		// 🔁 Replace this with your Kafka logic
		fmt.Printf("Publishing to Kafka topic: %s from job: %+v\n", job.KafkaTopic, job)
		scs.kafkaService.PublishJob(job.UserAlgorithmID.String(), CRON_JOB_EVENT_TYPE, job)
	}
}
