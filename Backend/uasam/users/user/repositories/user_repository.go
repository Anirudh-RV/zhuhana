package repositories

import (
	"database/sql"

	"uasam/middleware"
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

func (ur *UserRepository) CreateUser(firstname string, middleName *string, lastName string, emailID string, password string) (*models.User, error) {
	encFirstName, err := middleware.Encrypt(firstname)
	if err != nil {
		return nil, err
	}
	var encMiddleName *string
	if middleName != nil {
		e, err := middleware.Encrypt(*middleName)
		if err != nil {
			return nil, err
		}
		encMiddleName = &e
	}
	encLastName, err := middleware.Encrypt(lastName)
	if err != nil {
		return nil, err
	}
	encEmailID, err := middleware.EncryptDeterministic(emailID)
	if err != nil {
		return nil, err
	}

	// Insert into DB
	createQuery := `
		INSERT INTO "account" (first_name, middle_name, last_name, email_id, password)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`

	var user models.User
	err = ur.db.QueryRow(createQuery, encFirstName, encMiddleName, encLastName, encEmailID, password).
		Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}

	// Store plain values in model for app use
	user.FirstName = firstname
	user.MiddleName = middleName
	user.LastName = lastName
	user.EmailID = emailID
	user.Password = password

	return &user, nil
}

func (ur *UserRepository) IfUserEmailExists(emailID string) (bool, error) {
	encEmail, err := middleware.EncryptDeterministic(emailID)
	if err != nil {
		return false, err
	}
	query := `SELECT 1 FROM "account" WHERE email_id = $1 LIMIT 1`
	var exists int
	err = ur.db.QueryRow(query, encEmail).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

func (ur *UserRepository) GetUserPasswordByEmail(emailID string) (string, error) {
	encEmailID, err := middleware.EncryptDeterministic(emailID)
	if err != nil {
		return "", err
	}
	query := `SELECT password FROM "account" WHERE email_id = $1 LIMIT 1`
	var password string
	err = ur.db.QueryRow(query, encEmailID).Scan(&password)
	if err == sql.ErrNoRows {
		return "", nil
	}
	return password, err
}

func (ur *UserRepository) GetUserByEmail(emailID string) (*models.UserObject, error) {
	encEmailID, err := middleware.EncryptDeterministic(emailID)
	if err != nil {
		return nil, err
	}

	query := `
		SELECT id, first_name, middle_name, last_name, email_id, created_at, updated_at
		FROM "account"
		WHERE email_id = $1
		LIMIT 1
	`

	var encFirst, encLast, encEmailOut string
	var encMiddle *string
	var user models.UserObject

	err = ur.db.QueryRow(query, encEmailID).Scan(
		&user.ID,
		&encFirst,
		&encMiddle,
		&encLast,
		&encEmailOut,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	// Decrypt PII fields
	user.FirstName, _ = middleware.Decrypt(encFirst)
	user.LastName, _ = middleware.Decrypt(encLast)
	user.EmailID, _ = middleware.DecryptDeterministic(encEmailOut)
	if encMiddle != nil {
		d, _ := middleware.Decrypt(*encMiddle)
		user.MiddleName = &d
	}

	return &user, nil
}

func (ur *UserRepository) UpdateUserPasswordByEmail(emailID string, newPassword string) error {
	encEmailID, err := middleware.EncryptDeterministic(emailID)
	if err != nil {
		return err
	}
	query := `UPDATE "account" SET password = $1, updated_at = NOW() WHERE email_id = $2`
	_, err = ur.db.Exec(query, newPassword, encEmailID)
	return err
}
