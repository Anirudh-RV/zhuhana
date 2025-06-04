package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type OrderDomain int

const (
	Backtest OrderDomain = iota
	PaperTrading
	LiveTrading
)

func (o OrderDomain) String() string {
	switch o {
	case Backtest:
		return "Backtest"
	case PaperTrading:
		return "PaperTrading"
	case LiveTrading:
		return "LiveTrading"
	default:
		return "Unknown"
	}
}

func (o OrderDomain) Int() string {
	return o.Int()
}

type UserAlgorithm struct {
	ID           uuid.UUID   `db:"id"`
	UserID       uuid.UUID   `db:"user_id"`
	ScriptName   string      `db:"script_name"`
	ScriptURL    string      `db:"script_url"`
	CronSchedule string      `db:"cron_schedule"`
	OrderDomain  OrderDomain `db:"order_domain"`
	CreatedAt    time.Time   `db:"created_at"`
	UpdatedAt    time.Time   `db:"updated_at"`
}

type UserAlgorithmInfo struct {
	ScriptID     uuid.UUID   `json:"scriptID"`
	ScriptName   string      `json:"scriptName"`
	ScriptURL    *string     `json:"script_url,omitempty"`
	CronSchedule *string     `json:"cronSchedule,omitempty"`
	OrderDomain  OrderDomain `json:"order_domain"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
}

func (o OrderDomain) MarshalJSON() ([]byte, error) {
	return fmt.Appendf(nil, "\"%s\"", o.String()), nil
}

func (o *OrderDomain) UnmarshalJSON(data []byte) error {
	switch string(data) {
	case `"Backtest"`:
		*o = Backtest
	case `"PaperTrading"`:
		*o = PaperTrading
	case `"LiveTrading"`:
		*o = LiveTrading
	default:
		return fmt.Errorf("invalid OrderDomain: %s", data)
	}
	return nil
}
