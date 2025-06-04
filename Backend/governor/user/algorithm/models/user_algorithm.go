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

type UserAlgorithmInfo struct {
	ScriptID     uuid.UUID `json:"scriptID"`
	ScriptName   string    `json:"scriptName"`
	ScriptURL    *string   `json:"script_url,omitempty"`
	CronSchedule *string   `json:"cronSchedule,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
