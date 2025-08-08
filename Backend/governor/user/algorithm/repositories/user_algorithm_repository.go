package repositories

import (
	"database/sql"
	"fmt"
	"governor/user/algorithm/models"

	"github.com/google/uuid"
)

type UserAlgorithmRepository struct {
	db *sql.DB
}

func NewUserAlgorithmRepository(db *sql.DB) *UserAlgorithmRepository {
	return &UserAlgorithmRepository{
		db: db,
	}
}

func (uar *UserAlgorithmRepository) CreateUserAlgorithm(userID, scriptName string) (*models.UserAlgorithm, error) {

	// Insert into DB
	createQuery := `
		INSERT INTO "user_algorithm" (user_id, script_name)
		VALUES ($1, $2)
		RETURNING id, user_id, created_at, updated_at
	`

	var userAlgorithm models.UserAlgorithm
	err := uar.db.QueryRow(createQuery, userID, scriptName).
		Scan(&userAlgorithm.ID, &userAlgorithm.UserID, &userAlgorithm.CreatedAt, &userAlgorithm.UpdatedAt)
	if err != nil {
		return nil, err
	}

	// Store plain values in model for app use
	userAlgorithm.ScriptName = scriptName
	return &userAlgorithm, nil
}

func (uar *UserAlgorithmRepository) UpdateUserAlgorithmScriptName(userID string, algorithmID uuid.UUID, newScriptName string) (*models.UserAlgorithm, error) {
	updateQuery := `
		UPDATE "user_algorithm"
		SET script_name = $1, updated_at = NOW()
		WHERE id = $2 AND user_id = $3
		RETURNING id, user_id, script_name, script_url, start_cron_schedule, end_cron_schedule, order_domain, created_at, updated_at
	`

	var updatedAlgorithm models.UserAlgorithm
	err := uar.db.QueryRow(
		updateQuery,
		newScriptName,
		algorithmID,
		userID,
	).Scan(
		&updatedAlgorithm.ID,
		&updatedAlgorithm.UserID,
		&updatedAlgorithm.ScriptName,
		&updatedAlgorithm.ScriptURL,
		&updatedAlgorithm.StartCronSchedule,
		&updatedAlgorithm.EndCronSchedule,
		&updatedAlgorithm.OrderDomain,
		&updatedAlgorithm.CreatedAt,
		&updatedAlgorithm.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &updatedAlgorithm, nil
}

func (uar *UserAlgorithmRepository) UpdateScriptURL(userAlgorithmID, scriptURL string) error {
	query := `UPDATE "user_algorithm" SET script_url = $1, updated_at = NOW() WHERE id = $2`
	_, err := uar.db.Exec(query, scriptURL, userAlgorithmID)
	return err
}

func (uar *UserAlgorithmRepository) UpdateCronSchedule(userID, userAlgorithmID, startCronSchedule, endCronSchedule string) error {
	query := `UPDATE "user_algorithm" SET start_cron_schedule = $1, end_cron_schedule = $2, updated_at = NOW() WHERE id = $3 AND user_id = $4`
	res, err := uar.db.Exec(query, startCronSchedule, endCronSchedule, userAlgorithmID, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no matching user_algorithm found to update")
	}

	return err
}

func (uar *UserAlgorithmRepository) GetAllUserAlgorithmByUserID(userID string) ([]models.UserAlgorithmInfo, error) {
	query := `
		SELECT id, script_name, start_cron_schedule, end_cron_schedule, script_url, created_at, updated_at
		FROM "user_algorithm"
		WHERE user_id = $1
	`

	rows, err := uar.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var scripts []models.UserAlgorithmInfo
	for rows.Next() {
		var startCron sql.NullString
		var endCron sql.NullString
		var url sql.NullString
		var script models.UserAlgorithmInfo

		if err := rows.Scan(
			&script.ID,
			&script.ScriptName,
			&startCron,
			&endCron,
			&url,
			&script.CreatedAt,
			&script.UpdatedAt,
		); err != nil {
			return nil, err
		}

		if startCron.Valid {
			script.StartCronSchedule = &startCron.String
		} else {
			script.StartCronSchedule = nil
		}

		if endCron.Valid {
			script.EndCronSchedule = &endCron.String
		} else {
			script.EndCronSchedule = nil
		}

		if url.Valid {
			script.ScriptURL = &url.String
		} else {
			script.ScriptURL = nil
		}

		scripts = append(scripts, script)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return scripts, nil
}

func (uar *UserAlgorithmRepository) GetUserAlgorithmByUserID(userID, algorithmID string) (*models.UserAlgorithmInfo, error) {
	query := `
		SELECT id, script_name, start_cron_schedule, end_cron_schedule, script_url, created_at, updated_at
		FROM "user_algorithm"
		WHERE user_id = $1 AND id = $2
	`

	var userAlgorithm models.UserAlgorithmInfo
	var startCron sql.NullString
	var endCron sql.NullString
	var url sql.NullString

	err := uar.db.QueryRow(query, userID, algorithmID).Scan(
		&userAlgorithm.ID,
		&userAlgorithm.ScriptName,
		&startCron,
		&endCron,
		&url,
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

	return &userAlgorithm, nil
}

func (uar *UserAlgorithmRepository) DoesUserAlgorithmBelongsToUser(userID, userAlgorithmID string) (bool, error) {
	query := `SELECT 1 FROM "user_algorithm" WHERE id = $1 AND user_id = $2 LIMIT 1`
	rows, err := uar.db.Query(query, userAlgorithmID, userID)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	if rows.Next() {
		return true, nil // Belongs to user
	}

	return false, nil // No matching row found
}

func (uar *UserAlgorithmRepository) GetUserAlgorithmRunsByUserAlgorithmID(userAlgorithmID uuid.UUID) ([]models.UserAlgorithmRun, error) {
	query := `
		SELECT
			id,
			is_active,
			user_algorithm_id,
			start_cron_schedule,
			end_cron_schedule,
			order_domain,
			market,
			symbol,
			start_time,
			end_time,
			frequency,
			portfolio_size,
			status,
			created_at,
			stopped_at,
			updated_at
		FROM user_algorithm_runs
		WHERE user_algorithm_id = $1
		ORDER BY created_at DESC
	`

	rows, err := uar.db.Query(query, userAlgorithmID)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var runs []models.UserAlgorithmRun
	for rows.Next() {
		var run models.UserAlgorithmRun
		err := rows.Scan(
			&run.ID,
			&run.IsActive,
			&run.UserAlgorithmID,
			&run.StartCronSchedule,
			&run.EndCronSchedule,
			&run.OrderDomain,
			&run.Market,
			&run.Symbol,
			&run.StartTime,
			&run.EndTime,
			&run.Frequency,
			&run.PortfolioSize,
			&run.Status,
			&run.CreatedAt,
			&run.StoppedAt,
			&run.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}
		runs = append(runs, run)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return runs, nil
}
