package models

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)


// Rol es una estructura con 2 potenciales valores, ADMIN o REGULAR.

type Rol string

const (
	ADMIN   Rol = "admin"
	REGULAR Rol = "regular"
	ALL     Rol = "all"
)

var ADMIN_ROLES = []Rol{ADMIN}
var ALL_ROLES = []Rol{ADMIN, REGULAR}

// postgres

// en este caso si quisieramos dejar type Rol string directamente
// gorm tira error porque no sabe como mapearlo.
// se dejo como un array de postgres.

type UserPostgres struct {
	gorm.Model
	Name       string         `json:"name" bson:"name"`
	MiddleName string         `json:"middlename" bson:"middlename"`
	Email      string         `json:"email" bson:"email"`
	Password   string         `json:"password" bson:"password"`
	Phone      string         `json:"phone" bson:"phone"`
	Roles      pq.StringArray `gorm:"type:text[]" json:"roles" bson:"roles"`
	ActiveRol  string         `json:"active_rol" bson:"active_rol"`
}	
