package repositories

import (
	"database/sql"
	"fmt"
	"governor/user/algorithm/models"
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

func (uar *UserAlgorithmRepository) UpdateScriptURL(scriptID, scriptURL string) error {
	query := `UPDATE "user_algorithm" SET script_url = $1, updated_at = NOW() WHERE id = $2`
	_, err := uar.db.Exec(query, scriptURL, scriptID)
	return err
}

func (uar *UserAlgorithmRepository) UpdateCronSchedule(userID, scriptID, cronSchedule string) error {
	query := `UPDATE "user_algorithm" SET cron_schedule = $1, updated_at = NOW() WHERE id = $2 AND user_id = $3`
	res, err := uar.db.Exec(query, cronSchedule, scriptID, userID)
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
		SELECT id, script_name, cron_schedule, script_url, created_at, updated_at
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
		var cron sql.NullString
		var url sql.NullString
		var script models.UserAlgorithmInfo

		if err := rows.Scan(
			&script.ScriptID,
			&script.ScriptName,
			&cron,
			&url,
			&script.CreatedAt,
			&script.UpdatedAt,
		); err != nil {
			return nil, err
		}

		if cron.Valid {
			script.CronSchedule = &cron.String
		} else {
			script.CronSchedule = nil
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
		SELECT id, script_name, cron_schedule, script_url, created_at, updated_at
		FROM "user_algorithm"
		WHERE user_id = $1 AND id = $2
	`

	var userAlgorithm models.UserAlgorithmInfo
	var cron sql.NullString
	var url sql.NullString

	err := uar.db.QueryRow(query, userID, algorithmID).Scan(
		&userAlgorithm.ScriptID,
		&userAlgorithm.ScriptName,
		&cron,
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

	if cron.Valid {
		userAlgorithm.CronSchedule = &cron.String
	} else {
		userAlgorithm.CronSchedule = nil
	}

	if url.Valid {
		userAlgorithm.ScriptURL = &url.String
	} else {
		userAlgorithm.ScriptURL = nil
	}

	return &userAlgorithm, nil
}
