package commonutils

import (
	"errors"
	"fmt"
	"strings"

	"github.com/robfig/cron/v3"
)

type CronFields struct {
	Minute     string
	Hour       string
	DayOfMonth string
	Month      string
	DayOfWeek  string
}

// ValidateCronExpression uses robfig/cron to ensure it's syntactically correct.
func ValidateCronExpression(expr string) error {
	// robfig/cron expects standard 5-field or 6-field expressions (with seconds).
	// For 5-field, use cron.NewParser with appropriate options.
	parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	_, err := parser.Parse(expr)
	return err
}

// ParseCronExpression splits a valid cron expression into fields.
func ParseCronExpression(expr string) (*CronFields, error) {
	parts := strings.Fields(expr)
	if len(parts) != 5 {
		return nil, errors.New("invalid cron expression: must have 5 fields")
	}

	return &CronFields{
		Minute:     parts[0],
		Hour:       parts[1],
		DayOfMonth: parts[2],
		Month:      parts[3],
		DayOfWeek:  parts[4],
	}, nil
}

// ValidateAndParseCron validates and returns the split cron fields.
func ValidateAndParseCron(expr string) (*CronFields, error) {
	if err := ValidateCronExpression(expr); err != nil {
		return nil, fmt.Errorf("invalid cron expression: %w", err)
	}

	return ParseCronExpression(expr)
}
