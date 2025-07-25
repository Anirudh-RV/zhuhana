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
	upsertQuery := `
		INSERT INTO "user_secret" (user_id, key, value)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, key) DO UPDATE SET
			value = EXCLUDED.value,
			updated_at = CURRENT_TIMESTAMP; -- Optional: track updates
	`

	_, err = ur.db.Exec(upsertQuery, userUUID, encKey, encValue)
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

func (ur *UserSecretRepository) GetAllUserSecretKeys(userID string) ([]string, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	query := `
		SELECT key
		FROM "user_secret"
		WHERE user_id = $1
	`

	rows, err := ur.db.Query(query, userUUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var encKeys []string
	for rows.Next() {
		var encKey string
		if err := rows.Scan(&encKey); err != nil {
			return nil, err
		}
		encKeys = append(encKeys, encKey)
	}

	// Decrypt each key
	var decryptedKeys []string
	for _, ek := range encKeys {
		key, err := middleware.DecryptDeterministic(ek)
		if err != nil {
			continue
		}
		decryptedKeys = append(decryptedKeys, key)
	}

	return decryptedKeys, nil
}

func (ur *UserSecretRepository) DeleteUserSecretByID(userID, secretID string) error {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return err
	}
	secretUUID, err := uuid.Parse(secretID)
	if err != nil {
		return err
	}

	query := `
		DELETE FROM "user_secret"
		WHERE user_id = $1 AND id = $2
	`

	result, err := ur.db.Exec(query, userUUID, secretUUID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows // or a custom error like ErrSecretNotFound
	}

	return nil
}
