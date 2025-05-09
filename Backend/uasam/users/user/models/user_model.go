package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID         uuid.UUID `db:"id"`
	FirstName  string    `db:"first_name"`
	MiddleName *string   `db:"middle_name"`
	LastName   string    `db:"last_name"`
	EmailID    string    `db:"email_id"`
	Password   string    `db:"password"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}
