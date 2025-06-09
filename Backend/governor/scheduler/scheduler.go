package scheduler

import (
	"fmt"
	"log"

	"github.com/google/uuid"
)

func LoadCronJob() {
	jobs, err := GetAllActiveJobs()
	if err != nil {
		log.Fatalln("Failed to load jobs:", err)
	}

	// Step 3: Setup Cron Scheduler
	for _, job := range jobs {
		entryID, err := CronScheduler.AddFunc(job.Schedule, KafkaJobWrapper(job))
		if err := UpdateCronEntryID(job.ID, int64(entryID)); err != nil {
			log.Printf("Failed to schedule job %s: %v", job.UserAlgorithmID, err)
		}
		if err != nil {
			log.Printf("Failed to schedule job %s: %v", job.UserAlgorithmID, err)
		} else {
			fmt.Printf("Scheduled job %s with EntryID: %d\n", job.UserAlgorithmID, entryID)
		}
	}
}

func ScheduleCronJob(userAlgorithmID uuid.UUID, schedule, jobType, kafkaTopic string) error {
	// 1. Insert into DB
	jobID, err := InsertJob(userAlgorithmID, schedule, jobType, kafkaTopic)
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
	cronEntryID, err := CronScheduler.AddFunc(job.Schedule, KafkaJobWrapper(job))
	if err != nil {
		return fmt.Errorf("failed to add job to scheduler: %w", err)
	}
	job.CronEntryID = int64(cronEntryID)

	if err := UpdateCronEntryID(jobID, int64(cronEntryID)); err != nil {
		return fmt.Errorf("failed to update cron entry ID: %w", err)
	}

	// 3. Publish the Kafka Job
	log.Printf("✅ Scheduled and inserted job '%s' (%s)", userAlgorithmID, jobID)
	return nil
}
