package kubernetescontroller

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

type UserAlgorithmRun struct {
	ID                uuid.UUID   `db:"id"`
	IsActive          bool        `db:"is_active"`
	UserAlgorithmID   uuid.UUID   `db:"user_algorithm_id"`
	StartCronSchedule *string     `db:"start_cron_schedule"`
	EndCronSchedule   *string     `db:"end_cron_schedule"`
	OrderDomain       OrderDomain `db:"order_domain"`
	CreatedAt         time.Time   `db:"created_at"`
	StoppedAt         time.Time   `db:"stopped_at"`
	UpdatedAt         time.Time   `db:"updated_at"`
}
