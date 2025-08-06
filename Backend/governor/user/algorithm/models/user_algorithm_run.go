package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type RunStatus int

const (
	StatusCreating RunStatus = iota
	StatusRunning
	StatusCompleted
	StatusStopped
)

func (r RunStatus) String() string {
	switch r {
	case StatusCreating:
		return "Creating"
	case StatusRunning:
		return "Running"
	case StatusCompleted:
		return "Completed"
	case StatusStopped:
		return "Stopped"
	default:
		return "Unknown"
	}
}

func (r RunStatus) MarshalJSON() ([]byte, error) {
	return fmt.Appendf(nil, "\"%s\"", r.String()), nil
}

func (r *RunStatus) UnmarshalJSON(data []byte) error {
	switch string(data) {
	case `"Creating"`:
		*r = StatusCreating
	case `"Running"`:
		*r = StatusRunning
	case `"Completed"`:
		*r = StatusCompleted
	case `"Stopped"`:
		*r = StatusStopped
	default:
		return fmt.Errorf("invalid RunStatus: %s", data)
	}
	return nil
}

type UserAlgorithmRun struct {
	ID                uuid.UUID   `db:"id"`
	IsActive          bool        `db:"is_active"`
	UserAlgorithmID   uuid.UUID   `db:"user_algorithm_id"`
	StartCronSchedule *string     `db:"start_cron_schedule"`
	EndCronSchedule   *string     `db:"end_cron_schedule"`
	OrderDomain       OrderDomain `db:"order_domain"`
	Status            RunStatus   `db:"status"`
	Market            *string     `db:"market"`
	Symbol            *string     `db:"symbol"`
	StartTime         *time.Time  `db:"start_time"`
	EndTime           *time.Time  `db:"end_time"`
	Frequency         *int        `db:"frequency"`
	PortfolioSize     *int        `db:"portfolio_size"`
	CreatedAt         time.Time   `db:"created_at"`
	StoppedAt         *time.Time  `db:"stopped_at"`
	UpdatedAt         time.Time   `db:"updated_at"`
}

type UserAlgorithmLoginResponse struct {
	Status            int    `json:"status"`
	StatusDescription string `json:"statusDescription"`
	AccessToken       string `json:"accessToken"`
}
