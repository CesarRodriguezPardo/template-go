package services

import (
	"CesarRodriguezPardo/template-go/infra/database"
	"CesarRodriguezPardo/template-go/internal/repositories"
)

var (
	userRepo *repositories.UserRepository
)

func InitRepositories(db *database.Postgres) {
	userRepo = repositories.NewUserRepository(db)
}
