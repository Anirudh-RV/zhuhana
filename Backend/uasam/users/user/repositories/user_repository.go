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
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (ur *UserRepository) GetUserPasswordByEmail(emailID string) (string, error) {
	query := `SELECT password FROM "account" WHERE email_id = $1 LIMIT 1`
	var password string
	err := ur.db.QueryRow(query, emailID).Scan(&password)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", err
	}
	return password, nil
}

func (ur *UserRepository) GetUserByEmail(emailID string) (*models.UserObject, error) {
	query := `
		SELECT id, first_name, middle_name, last_name, email_id, created_at, updated_at
		FROM "account"
		WHERE email_id = $1
		LIMIT 1
	`

	var user models.UserObject
	err := ur.db.QueryRow(query, emailID).Scan(
		&user.ID,
		&user.FirstName,
		&user.MiddleName,
		&user.LastName,
		&user.EmailID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (ur *UserRepository) UpdateUserPasswordByEmail(emailID string, newPassword string) error {
	query := `UPDATE "account" SET password = $1, updated_at = NOW() WHERE email_id = $2`
	_, err := ur.db.Exec(query, newPassword, emailID)
	if err != nil {
		return err
	}
	return nil
}
