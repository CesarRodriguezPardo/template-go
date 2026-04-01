package repositories

import (
	"CesarRodriguezPardo/template-go/infra/database"
	"CesarRodriguezPardo/template-go/internal/models"
	"context"
	"fmt"

	uuid "github.com/satori/go.uuid"
)

type UserRepository struct {
	DB *database.Postgres
}

func (r *UserRepository) InsertOne(ctx context.Context, user *models.User) (uuid.UUID, error) {
	query := `
		INSERT INTO users (name, email, created_at)
		VALUES ($1, $2, NOW())
		RETURNING id
	`

	var id uuid.UUID
	err := r.DB.Pool().QueryRow(ctx, query,
		user.Name,
		user.Email,
	).Scan(&id)

	if err != nil {
		return uuid.UUID{}, fmt.Errorf("UserRepository.InsertOne: %w", err)
	}

	return id, nil
}

func NewUserRepository(db *database.Postgres) *UserRepository {
	return &UserRepository{
		DB: db,
	}
}
