package dto

// LoginRequest represents the credentials needed for authentication.
type LoginRequest struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required"`
}
