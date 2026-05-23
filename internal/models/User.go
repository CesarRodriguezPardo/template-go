package models

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	uuid "github.com/satori/go.uuid"
)

type User struct {
	ID         uuid.UUID `db:"id"`
	Name       string    `db:"name"`
	MiddleName string    `db:"middle_name"`
	Email      string    `db:"email"`
	Phone      string    `db:"phone"`

	Password string `db:"password"`

	Role string `db:"role"`

	CreatedAt time.Time          `db:"created_at"`
	UpdatedAt time.Time          `db:"updated_at"`
	DeletedAt pgtype.Timestamptz `db:"deleted_at"`
}
