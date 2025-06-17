package scheduler

import (
	"fmt"

	"github.com/google/uuid"
)

func (scs *SchedulerService) GetAllActiveJobs() ([]CronJob, error) {
	query := `
		SELECT id, user_algorithm_id, schedule, job_type, kafka_topic, is_active, created_at, updated_at
		FROM cron_job
		WHERE is_active = true
	`

	rows, err := scs.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []CronJob
	for rows.Next() {
		var job CronJob
		if err := rows.Scan(
			&job.ID,
			&job.UserAlgorithmID,
			&job.Schedule,
			&job.JobType,
			&job.KafkaTopic,
			&job.IsActive,
			&job.CreatedAt,
			&job.UpdatedAt,
		); err != nil {
			return nil, err
		}
		jobs = append(jobs, job)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return jobs, nil
}

func (scs *SchedulerService) GetAllJobsForUserAlgorithmWithJobType(userAlgorithmID uuid.UUID, jobType string) ([]int64, error) {
	query := `
		SELECT cron_entry_id
		FROM cron_job
		WHERE is_active = true AND user_algorithm_id = $1 AND job_type = $2
	`

	rows, err := scs.db.Query(query, userAlgorithmID, jobType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cronEntries []int64
	for rows.Next() {
		var cronEntryID int64
		if err := rows.Scan(
			&cronEntryID,
		); err != nil {
			return nil, err
		}
		cronEntries = append(cronEntries, cronEntryID)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return cronEntries, nil
}

func (scs *SchedulerService) InsertJob(userAlgorithmID uuid.UUID, schedule, jobType, kafkaTopic string) (uuid.UUID, error) {
	query := `
		INSERT INTO cron_job (user_algorithm_id, schedule, job_type, kafka_topic)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	var id uuid.UUID
	err := scs.db.QueryRow(query, userAlgorithmID, schedule, jobType, kafkaTopic).Scan(&id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("insert failed: %w", err)
	}
	return id, nil
}

func (scs *SchedulerService) DeactivateUserAlgorithmWithJobType(userAlgorithmID uuid.UUID, jobType string) error {
	query := `
		UPDATE cron_job
		SET is_active = false
		WHERE user_algorithm_id = $1 AND job_type = $2
	`

	_, err := scs.db.Exec(query, userAlgorithmID, jobType)
	return err
}

func (scs *SchedulerService) DeactivateJob(id uuid.UUID) error {
	query := `
		UPDATE cron_job
		SET is_active = false
		WHERE id = $1
	`

	_, err := scs.db.Exec(query, id)
	return err
}

func (scs *SchedulerService) UpdateCronEntryID(jobID uuid.UUID, entryID int64) error {
	query := `
		UPDATE cron_job
		SET cron_entry_id = $1
		WHERE id = $2
	`
	_, err := scs.db.Exec(query, entryID, jobID)
	return err
}

func (scs *SchedulerService) DeactivateUserAlgorithm(userAlgorithmID uuid.UUID) error {
	query := `
		UPDATE cron_job
		SET is_active = false
		WHERE user_algorithm_id = $1
	`

	_, err := scs.db.Exec(query, userAlgorithmID)
	return err
}

func (scs *SchedulerService) GetAllJobsForUserAlgorithm(userAlgorithmID uuid.UUID) ([]int64, error) {
	query := `
		SELECT cron_entry_id
		FROM cron_job
		WHERE is_active = true AND user_algorithm_id = $1
	`

	rows, err := scs.db.Query(query, userAlgorithmID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cronEntries []int64
	for rows.Next() {
		var cronEntryID int64
		if err := rows.Scan(
			&cronEntryID,
		); err != nil {
			return nil, err
		}
		cronEntries = append(cronEntries, cronEntryID)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return cronEntries, nil
}
