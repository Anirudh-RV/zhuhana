package scheduler

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/robfig/cron/v3"
)

func (scs *SchedulerService) Init() {
	scs.LoadCronJob()
	scs.StartRedisSubscriber(context.Background())
}

func (scs *SchedulerService) LoadCronJob() {
	jobs, err := scs.GetAllActiveJobs()
	if err != nil {
		log.Fatalln("Failed to load jobs:", err)
	}

	for _, job := range jobs {
		ctx := context.Background()
		lockKey := fmt.Sprintf("cron-schedule-lock:%s", job.ID)

		// Only one instance should schedule this cron
		lock, err := scs.TryLock(ctx, lockKey, 2*time.Minute) // use a slightly longer lock for scheduling phase
		if err != nil {
			fmt.Println("another node is scheduling this job:", job.UserAlgorithmID)
			continue
		}

		entryID, err := scs.cronScheduler.AddFunc(job.Schedule, scs.KafkaJobWrapper(job))
		if err != nil {
			log.Printf("Failed to schedule job %s: %v", job.UserAlgorithmID, err)
			lock.Release(ctx) // unlock if scheduling failed
			continue
		}

		if err := scs.UpdateCronEntryID(job.ID, int64(entryID)); err != nil {
			log.Printf("Failed to update entry ID for job %s: %v", job.UserAlgorithmID, err)
		} else {
			fmt.Printf("Scheduled job %s with EntryID: %d\n", job.UserAlgorithmID, entryID)
		}

		lock.Release(ctx)
	}
}

func (scs *SchedulerService) ScheduleCronJob(userAlgorithmID uuid.UUID, schedule, jobType, kafkaTopic string) error {
	// 1. Insert into DB
	if err := scs.CancelCronJobForUserAlgorithmWithJobType(userAlgorithmID, jobType); err != nil {
		return fmt.Errorf("failed to deactivate other user algorithms: %w", err)
	}

	if err := scs.DeactivateUserAlgorithmWithJobType(userAlgorithmID, jobType); err != nil {
		return fmt.Errorf("failed to deactivate other user algorithms: %w", err)
	}

	jobID, err := scs.InsertJob(userAlgorithmID, schedule, jobType, kafkaTopic)
	if err != nil {
		return fmt.Errorf("failed to insert cron job: %w", err)
	}

	job := CronJob{
		ID:              jobID,
		UserAlgorithmID: userAlgorithmID,
		Schedule:        schedule,
		JobType:         jobType,
		KafkaTopic:      kafkaTopic,
		IsActive:        true,
	}

	// 2. Add to scheduler with Redis lock
	cronEntryID, err := scs.cronScheduler.AddFunc(job.Schedule, scs.KafkaJobWrapper(job))
	if err != nil {
		return fmt.Errorf("failed to add job to scheduler: %w", err)
	}
	job.CronEntryID = int64(cronEntryID)

	if err := scs.UpdateCronEntryID(jobID, int64(cronEntryID)); err != nil {
		return fmt.Errorf("failed to update cron entry ID: %w", err)
	}

	// 3. Publish the Kafka Job
	log.Printf("✅ Scheduled and inserted job '%s' (%s)", userAlgorithmID, jobID)
	return nil
}

func (scs *SchedulerService) CancelCronJobForUserAlgorithmWithJobType(userAlgorithmID uuid.UUID, jobType string) error {
	cronEntries, err := scs.GetAllJobsForUserAlgorithmWithJobType(userAlgorithmID, jobType)
	if err != nil {
		return err
	}
	for _, cronEntryID := range cronEntries {
		scs.CancelCronJobWithID(userAlgorithmID, cronEntryID)
	}

	return nil
}

func (scs *SchedulerService) CancelCronJobWithID(userAlgorithmID uuid.UUID, entryID int64) {
	scs.cronScheduler.Remove(cron.EntryID(entryID))
	log.Printf("❌ Cron job with EntryID %d has been cancelled locally", entryID)
	log.Printf("Broadcasting cancellation of cron job with EntryID %d", entryID)
	_ = scs.BroadcastCancelJob(userAlgorithmID, entryID)
}

func (scs *SchedulerService) CancelCronJobForUserAlgorithm(userAlgorithmID uuid.UUID) error {
	cronEntries, err := scs.GetAllJobsForUserAlgorithm(userAlgorithmID)
	if err != nil {
		return err
	}
	for _, cronEntryID := range cronEntries {
		scs.CancelCronJobWithID(userAlgorithmID, cronEntryID)
	}

	scs.DeactivateUserAlgorithm(userAlgorithmID)
	return nil
}
