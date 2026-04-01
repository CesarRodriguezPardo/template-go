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

func (repo *UserRepository) CreateUser(ctx context.Context, user *models.User) (uuid.UUID, error) {
	query := `
		INSERT INTO users (name, middle_name, email, phone, password, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
		RETURNING id
	`

	var id uuid.UUID
	err := repo.DB.Pool().QueryRow(ctx, query,
		user.Name,
		user.MiddleName,
		user.Email,
		user.Phone,
		user.Password,
		user.Role,
	).Scan(&id)

	if err != nil {
		return uuid.UUID{}, fmt.Errorf("UserRepository.InsertOne: %w", err)
	}

	return id, nil
}

func (repo *UserRepository) FindUserByEmail(ctx context.Context, email string) (uuid.UUID, error) {
	query := `
		SELECT id FROM users
		WHERE email := $1
	`
	var id uuid.UUID
	err := repo.DB.Pool().QueryRow(ctx, query, email).Scan(&id)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("UserRepository.findUserByEmail: %w", err)
	}

	return id, nil
}

func NewUserRepository(db *database.Postgres) *UserRepository {
	return &UserRepository{
		DB: db,
	}
}
