package models

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	uuid "github.com/satori/go.uuid"
)

type User struct {
	ID         uuid.UUID `json:"id"         db:"id"`
	Name       string    `json:"name"       db:"name"`
	MiddleName string    `json:"middlename" db:"middle_name"`
	Email      string    `json:"email"      db:"email"`
	Phone      string    `json:"phone"      db:"phone"`

	Password string `json:"-" db:"password"`

	Role string `json:"role" db:"role"`

	CreatedAt time.Time          `json:"created_at" db:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" db:"updated_at"`
	DeletedAt pgtype.Timestamptz `json:"deleted_at" db:"deleted_at"`
}
