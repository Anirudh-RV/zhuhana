package scheduler

import (
	"time"

	"github.com/google/uuid"
)

var START_USER_ALGORITHM_JOB = "start-user-algorithm-job"
var END_USER_ALGORITHM_JOB = "end-user-algorithm-job"

type CronJob struct {
	ID              uuid.UUID `db:"id"`
	UserAlgorithmID uuid.UUID `db:"user_algorithm_id"`
	CronEntryID     int64     `db:"cron_entry_id"`
	Schedule        string    `db:"schedule"`
	JobType         string    `db:"job_type"`
	KafkaTopic      string    `db:"kafka_topic"`
	IsActive        bool      `db:"is_active"`
	CreatedAt       time.Time `db:"created_at"`
	UpdatedAt       time.Time `db:"updated_at"`
}
