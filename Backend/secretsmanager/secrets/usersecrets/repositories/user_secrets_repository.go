package repositories

import (
	"database/sql"
	"fmt"

	"secretsmanager/middleware"
	"secretsmanager/secrets/usersecrets/models"

	"github.com/google/uuid"
)

type UserSecretRepository struct {
	db *sql.DB
}

func NewUserSecretRepository(db *sql.DB) *UserSecretRepository {
	return &UserSecretRepository{
		db: db,
	}
}

func (ur *UserSecretRepository) CreateUserSecret(userID, key, value string) error {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return err
	}

	fmt.Println("Parsed UUID:", userUUID)

	encKey, err := middleware.EncryptDeterministic(key)
	if err != nil {
		return err
	}
	encValue, err := middleware.Encrypt(value)
	if err != nil {
		return err
	}

	// Insert into DB
	createQuery := `
		INSERT INTO "user_secret" (user_id, key, value)
		VALUES ($1, $2, $3)
	`

	_, err = ur.db.Exec(createQuery, userUUID, encKey, encValue)
	if err != nil {
		return err
	}

	return nil
}

func (ur *UserSecretRepository) GetUserSecret(userID, key string) (*models.UserSecret, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}
	encKey, err := middleware.EncryptDeterministic(key)
	if err != nil {
		return nil, err
	}

	query := `
		SELECT id, user_id, key, value, created_at, updated_at
		FROM "user_secret"
		WHERE user_id = $1 AND key = $2
		LIMIT 1
	`

	var encValue string
	var userSecret models.UserSecret
	userSecret.Key = key

	err = ur.db.QueryRow(query, userUUID, encKey).Scan(
		&userSecret.ID,
		&userSecret.UserID,
		&encKey,
		&encValue,
		&userSecret.CreatedAt,
		&userSecret.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	// Decrypt Value field
	userSecret.Value, _ = middleware.Decrypt(encValue)

	return &userSecret, nil
}
