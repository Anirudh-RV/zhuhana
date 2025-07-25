package repositories

import (
	"database/sql"
	"encoding/json"
	"uasam/users/user/models"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type NotificationRepository struct {
	db *sql.DB
}

func NewNotificationRepository(db *sql.DB) *NotificationRepository {
	return &NotificationRepository{
		db: db,
	}
}

func (nr *NotificationRepository) CreateNotification(
	userID uuid.UUID,
	notificationType, title, message string,
	link *string,
	metadata map[string]interface{},
) (*models.Notification, error) {

	// Serialize metadata to JSON
	metadataBytes, err := json.Marshal(metadata)
	if err != nil {
		return nil, err
	}

	query := `
		INSERT INTO notification (user_id, type, title, message, link, metadata)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, read, pinned, created_at, updated_at
	`

	var notif models.Notification
	err = nr.db.QueryRow(
		query,
		userID,
		notificationType,
		title,
		message,
		link,
		metadataBytes,
	).Scan(
		&notif.ID,
		&notif.Read,
		&notif.Pinned,
		&notif.CreatedAt,
		&notif.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	notif.UserID = userID
	notif.Type = notificationType
	notif.Title = title
	notif.Message = message
	notif.Link = link
	notif.Metadata.Bytes = metadataBytes

	return &notif, nil
}

func (nr *NotificationRepository) MarkNotificationsAsRead(ids []uuid.UUID) error {
	if len(ids) == 0 {
		return nil // no-op
	}

	query := `
		UPDATE notification
		SET read = TRUE, updated_at = NOW()
		WHERE id = ANY($1)
	`

	_, err := nr.db.Exec(query, pq.Array(ids))
	return err
}

func (nr *NotificationRepository) GetNotificationsByUserID(userID uuid.UUID) ([]models.Notification, error) {
	query := `
		SELECT id, user_id, type, title, message, link, read, pinned, metadata, created_at, updated_at
		FROM notification
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := nr.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []models.Notification

	for rows.Next() {
		var notif models.Notification
		err := rows.Scan(
			&notif.ID,
			&notif.UserID,
			&notif.Type,
			&notif.Title,
			&notif.Message,
			&notif.Link,
			&notif.Read,
			&notif.Pinned,
			&notif.Metadata,
			&notif.CreatedAt,
			&notif.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		notifications = append(notifications, notif)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return notifications, nil
}
