package repositories

import (
	"database/sql"

	"uasam/users/user/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

type UserRepositoryInterface interface {
	CreateUser(firstname string, middleName *string, lastName string, emailID string, password string) (*models.User, error)
	IfUserEmailExists(emailID string) (bool, error)
}

func (ur *UserRepository) CreateUser(firstname string, middleName *string, lastName string, emailID string, password string) (*models.User, error) {
	createQuery := `
		INSERT INTO "account" (first_name, middle_name, last_name, email_id, password)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`

	var user models.User
	user.FirstName = firstname
	user.MiddleName = middleName
	user.LastName = lastName
	user.EmailID = emailID
	user.Password = password

	err := ur.db.QueryRow(createQuery, firstname, middleName, lastName, emailID, password).
		Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (ur *UserRepository) IfUserEmailExists(emailID string) (bool, error) {
	query := `SELECT 1 FROM "account" WHERE email_id = $1 LIMIT 1`
	var exists int
	err := ur.db.QueryRow(query, emailID).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil // account does not exist
		}
		return false, err // other DB error
	}
	return true, nil // account exists
}
