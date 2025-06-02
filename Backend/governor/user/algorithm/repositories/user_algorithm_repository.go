package repositories

import (
	"database/sql"
	"governor/user/algorithm/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (ur *UserRepository) CreateUser(userID, scriptName, scriptURL, cronSchedule string) (*models.UserAlgorithm, error) {

	// Insert into DB
	createQuery := `
		INSERT INTO "user_algorithm" (user_id, script_name, script_url, cron_schedule)
		VALUES ($1, $2, $3, $4)
		RETURNING id, user_id, created_at, updated_at
	`

	var userAlgorithm models.UserAlgorithm
	err := ur.db.QueryRow(createQuery, userID, scriptName, scriptURL).
		Scan(&userAlgorithm.ID, &userAlgorithm.UserID, &userAlgorithm.CreatedAt, &userAlgorithm.UpdatedAt)
	if err != nil {
		return nil, err
	}

	// Store plain values in model for app use
	userAlgorithm.ScriptName = scriptName
	userAlgorithm.ScriptURL = scriptURL
	userAlgorithm.CronSchedule = cronSchedule

	return &userAlgorithm, nil
}
