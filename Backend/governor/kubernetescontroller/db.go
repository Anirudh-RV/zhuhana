package kubernetescontroller

import (
	"database/sql"
	"fmt"
	"governor/user/algorithm/models"

	"github.com/google/uuid"
)

func (ks *KubernetesService) GetUserAlgorithm(algorithmID uuid.UUID) (*models.UserAlgorithm, error) {
	query := `
		SELECT user_id, script_name, script_url, start_cron_schedule, end_cron_schedule, order_domain, created_at, updated_at
		FROM "user_algorithm"
		WHERE id = $1
	`

	var userAlgorithm models.UserAlgorithm
	var startCron sql.NullString
	var endCron sql.NullString
	var url sql.NullString

	err := ks.db.QueryRow(query, algorithmID).Scan(
		&userAlgorithm.UserID,
		&userAlgorithm.ScriptName,
		&userAlgorithm.ScriptURL,
		&startCron,
		&endCron,
		&userAlgorithm.OrderDomain,
		&userAlgorithm.CreatedAt,
		&userAlgorithm.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	if startCron.Valid {
		userAlgorithm.StartCronSchedule = &startCron.String
	} else {
		userAlgorithm.StartCronSchedule = nil
	}

	if endCron.Valid {
		userAlgorithm.EndCronSchedule = &endCron.String
	} else {
		userAlgorithm.EndCronSchedule = nil
	}
	if url.Valid {
		userAlgorithm.ScriptURL = &url.String
	} else {
		userAlgorithm.ScriptURL = nil
	}

	userAlgorithm.ID = algorithmID

	return &userAlgorithm, nil
}

func (ks *KubernetesService) AddUserAlgorithmRun(userAlgorithmID uuid.UUID, start_cron_schedule, end_cron_schedule string, order_domain int) (uuid.UUID, error) {
	query := `
		INSERT INTO user_algorithm_run (user_algorithm_id, start_cron_schedule, end_cron_schedule, order_domain)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	var id uuid.UUID
	err := ks.db.QueryRow(query, userAlgorithmID, start_cron_schedule, end_cron_schedule, order_domain).Scan(&id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("insert failed: %w", err)
	}
	return id, nil
}

func (ks *KubernetesService) GetUserAlgorithmRunByUserAlgorithmID(userAlgorithmID string) (uuid.UUID, error) {
	query := `
		SELECT id
		FROM "user_algorithm_run"
		WHERE user_algorithm_id = $1 AND is_active = true
	`

	var id uuid.UUID

	err := ks.db.QueryRow(query, userAlgorithmID).Scan(
		&id,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return uuid.Nil, nil
		}
		return uuid.Nil, err
	}

	return id, nil
}
