package dto

import (
	"time"

	"CesarRodriguezPardo/template-go/internal/models"
	uuid "github.com/satori/go.uuid"
)

// CreateUserRequest represents the input data for creating or updating a user.
type CreateUserRequest struct {
	Name       string `json:"name"       binding:"required"`
	MiddleName string `json:"middlename" binding:"required"`
	Email      string `json:"email"      binding:"required,email"`
	Password   string `json:"password"   binding:"required"`
	Phone      string `json:"phone"      binding:"required"`
	Role       string `json:"role"`
}

// UserResponse represents the output data returned to clients (excludes password).
type UserResponse struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	MiddleName string    `json:"middlename"`
	Email      string    `json:"email"`
	Phone      string    `json:"phone"`
	Role       string    `json:"role"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// ToModel converts a CreateUserRequest into a models.User.
func (req *CreateUserRequest) ToModel() *models.User {
	return &models.User{
		Name:       req.Name,
		MiddleName: req.MiddleName,
		Email:      req.Email,
		Password:   req.Password,
		Phone:      req.Phone,
		Role:       req.Role,
	}
}

// UserToResponse converts a models.User into a UserResponse.
func UserToResponse(user *models.User) *UserResponse {
	return &UserResponse{
		ID:         user.ID,
		Name:       user.Name,
		MiddleName: user.MiddleName,
		Email:      user.Email,
		Phone:      user.Phone,
		Role:       user.Role,
		CreatedAt:  user.CreatedAt,
		UpdatedAt:  user.UpdatedAt,
	}
}

// UsersToResponseList converts a slice of models.User to a slice of UserResponse.
func UsersToResponseList(users []*models.User) []*UserResponse {
	responses := make([]*UserResponse, len(users))
	for i, u := range users {
		responses[i] = UserToResponse(u)
	}
	return responses
}
