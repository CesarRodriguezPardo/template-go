package forms

import (
	"CesarRodriguezPardo/template-go/internal/models"
)

type UserForm struct {
	Name       string `json:"name"`
	MiddleName string `json:"middlename"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	Phone      string `json:"phone"`
	Role       string `json:"roles"`
}

func UserFormToUser(userForm *UserForm) *models.User {
	user := &models.User{}

	user.Name = userForm.Name
	user.MiddleName = userForm.MiddleName
	user.Email = userForm.Email
	user.Password = userForm.Password
	user.Phone = userForm.Phone
	user.Role = userForm.Role

	return user
}
