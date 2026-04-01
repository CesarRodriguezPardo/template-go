package forms

type UserForm struct {
	Name       string `json:"name"`
	MiddleName string `json:"middlename"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	Phone      string `json:"phone"`
	Role       string `json:"roles"`
}
