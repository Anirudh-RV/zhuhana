package repositories

import (
	"database/sql"
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

func (uar *UserAlgorithmRepository) UpdateCronSchedule(scriptID, cronSchedule string) error {
	query := `UPDATE "user_algorithm" SET cron_schedule = $1, updated_at = NOW() WHERE id = $2`
	_, err := uar.db.Exec(query, cronSchedule, scriptID)
	return err
}
