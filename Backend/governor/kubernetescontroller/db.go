package kubernetescontroller

import (
	"database/sql"
	"fmt"
	"governor/user/algorithm/models"
	"time"

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

func (ks *KubernetesService) AddUserAlgorithmRun(
	userAlgorithmID uuid.UUID,
	startCronSchedule,
	endCronSchedule string,
	orderDomain int,
	market,
	symbol string,
	startTime,
	endTime *time.Time,
	portfolioSize int,
) (uuid.UUID, error) {

	query := `
		INSERT INTO user_algorithm_runs (
			user_algorithm_id,
			start_cron_schedule,
			end_cron_schedule,
			order_domain,
			market,
			symbol,
			start_time,
			end_time,
			portfolio_size,
			created_at,
			updated_at,
			is_active
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW(), TRUE)
		RETURNING id
	`

	var id uuid.UUID
	err := ks.db.QueryRow(
		query,
		userAlgorithmID,
		startCronSchedule,
		endCronSchedule,
		orderDomain,
		market,
		symbol,
		startTime,
		endTime,
		portfolioSize,
	).Scan(&id)

	if err != nil {
		return uuid.Nil, fmt.Errorf("insert failed: %w", err)
	}

	return id, nil
}

func (ks *KubernetesService) GetUserAlgorithmRunsByUserAlgorithmID(userAlgorithmID string) ([]uuid.UUID, error) {
	query := `
		SELECT id
		FROM "user_algorithm_runs"
		WHERE user_algorithm_id = $1 AND is_active = true
	`

	rows, err := ks.db.Query(query, userAlgorithmID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []uuid.UUID
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return ids, nil
}

func (ks *KubernetesService) DeactivateUserAlgorithmRunByID(runID uuid.UUID) error {
	query := `
		UPDATE "user_algorithm_runs"
		SET is_active = false
		WHERE id = $1
	`

	result, err := ks.db.Exec(query, runID)
	if err != nil {
		return fmt.Errorf("failed to deactivate user_algorithm_runs %s: %w", runID, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("could not determine rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no user_algorithm_runs found with id %s", runID)
	}

	return nil
}
