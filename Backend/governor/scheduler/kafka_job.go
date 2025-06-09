package scheduler

import (
	"context"
	"fmt"
	"governor/kafka"
	"time"
)

func KafkaJobWrapper(job CronJob) func() {
	return func() {
		ctx := context.Background()
		lockKey := fmt.Sprintf("cron-lock:%s", job.ID)
		lock, err := TryLock(ctx, lockKey, time.Minute)
		if err != nil {
			fmt.Println("another node is handling this job:", job.UserAlgorithmID)
			return
		}
		defer lock.Release(ctx)

		// 🔁 Replace this with your Kafka logic
		fmt.Printf("Publishing to Kafka topic: %s from job: %s\n", job.KafkaTopic, job.UserAlgorithmID)
		kafka.PublishJob(job.UserAlgorithmID.String(), job)
	}
}
