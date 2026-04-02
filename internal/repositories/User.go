package repositories

import (
	"CesarRodriguezPardo/template-go/infra/database"
	"CesarRodriguezPardo/template-go/internal/models"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
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

func (repo *UserRepository) GetIdByEmail(ctx context.Context, email string) (uuid.UUID, error) {
	query := `
		SELECT id FROM users
		WHERE email = $1
	`
	var id uuid.UUID
	err := repo.DB.Pool().QueryRow(ctx, query, email).Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return uuid.Nil, nil
		}
		return uuid.Nil, fmt.Errorf("UserRepository.findUserByEmail: %w", err)
	}

	return id, nil
}

func (repo *UserRepository) GetAuthDataByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT password, role FROM users
		WHERE email = $1 
	`
	user := &models.User{}
	err := repo.DB.Pool().QueryRow(ctx, query, email).
		Scan(&user.Password, &user.Role)

	if err != nil {
		return nil, fmt.Errorf("credentials: %w", err)
	}

	return user, nil
}

/*
func (repo *UserRepository) GetAllUsers(ctx context.Context) ([]*models.User, error) {
	query := `
		SELECT name, middlename, email, password, phone, role FROM users
	`

	users := []*models.User{}




}
	*/

func NewUserRepository(db *database.Postgres) *UserRepository {
	return &UserRepository{
		DB: db,
	}
}
