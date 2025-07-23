package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgtype"
)

type Notification struct {
	ID        uuid.UUID    `db:"id"`
	UserID    uuid.UUID    `db:"user_id"`
	Type      string       `db:"type"`
	Title     string       `db:"title"`
	Message   string       `db:"message"`
	Link      *string      `db:"link"` // nullable
	Read      bool         `db:"read"`
	Pinned    bool         `db:"pinned"`
	Metadata  pgtype.JSONB `db:"metadata"` // or map[string]interface{} if you prefer
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt time.Time    `db:"updated_at"`
}

type GetNotificationsResponse struct {
	Status            int             `json:"status"`
	StatusDescription string          `json:"statusDescription"`
	Notifications     *[]Notification `json:"notifications,omitempty"`
}

type ReadNotificationsRequest struct {
	IDs []uuid.UUID `json:"IDs"`
}

type ReadNotificationsResponse struct {
	Status            int    `json:"status"`
	StatusDescription string `json:"statusDescription"`
}

type NotificationSwagger struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Type      string    `json:"type"`
	Title     string    `json:"title"`
	Message   string    `json:"message"`
	Link      *string   `json:"link,omitempty"`
	Read      bool      `json:"read"`
	Pinned    bool      `json:"pinned"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type GetNotificationsResponseSwagger struct {
	Status            int                    `json:"status"`
	StatusDescription string                 `json:"statusDescription"`
	Notifications     *[]NotificationSwagger `json:"notifications,omitempty"`
}
