package forms

import "citiaps/golang-backend-template/models"


// postgres
type UserFormPostgres struct {
	Name       string   `json:"name" bson:"name"`
	MiddleName string   `json:"middlename" bson:"middlename"`
	Email      string   `json:"email" bson:"email"`
	Password   string   `json:"password" bson:"password"`
	Phone      string   `json:"phone" bson:"phone"`
	Roles      []string `gorm:"type:json" json:"roles" bson:"roles"`
	ActiveRol  string   `json:"active_rol" bson:"active_rol"`
}
