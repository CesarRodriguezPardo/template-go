package models

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/lib/pq"
)

type User struct {
	Id         pgtype.UUID `json:"id"`
	Name       string `json:"name"`
	MiddleName string `json:"middlename"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`

	Password string `json:"password"`

	Role pq.StringArray `gorm:"type:text[]" json:"roles"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
