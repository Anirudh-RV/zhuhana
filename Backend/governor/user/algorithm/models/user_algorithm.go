package models

import (
	"time"

	"github.com/google/uuid"
)

type UserAlgorithm struct {
	ID           uuid.UUID `db:"id"`
	UserID       uuid.UUID `db:"user_id"`
	ScriptName   string    `db:"script_name"`
	ScriptURL    string    `db:"script_url"`
	CronSchedule string    `db:"cron_schedule"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}
