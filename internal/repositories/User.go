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
		WHERE email = $1 AND deleted_at IS NULL
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

func (repo *UserRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	query := `
		SELECT id, name, middle_name, email, phone, role, created_at, updated_at 
		FROM users 
		WHERE id = $1 AND deleted_at IS NULL
	`
	user := &models.User{}
	err := repo.DB.Pool().QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Name,
		&user.MiddleName,
		&user.Email,
		&user.Phone,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("UserRepository.GetUserByID: %w", err)
	}

	return user, nil
}

func (repo *UserRepository) GetAuthDataByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT id, password, role FROM users
		WHERE email = $1 AND deleted_at IS NULL
	`
	user := &models.User{}
	err := repo.DB.Pool().QueryRow(ctx, query, email).
		Scan(&user.ID, &user.Password, &user.Role)

	if err != nil {
		return nil, fmt.Errorf("credentials: %w", err)
	}

	return user, nil
}

func (repo *UserRepository) GetAllUsers(ctx context.Context, limit, offset int) ([]*models.User, int, error) {
	countQuery := `
		SELECT COUNT(*) FROM users WHERE deleted_at IS NULL
	`
	var total int
	if err := repo.DB.Pool().QueryRow(ctx, countQuery).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `
		SELECT id, name, middle_name, email, phone, role 
		FROM users
		WHERE deleted_at IS NULL
		LIMIT $1 OFFSET $2
	`

	rows, err := repo.DB.Pool().Query(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	users := []*models.User{}

	for rows.Next() {
		user := &models.User{}

		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.MiddleName,
			&user.Email,
			&user.Phone,
			&user.Role,
		)
		if err != nil {
			return nil, 0, err
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (repo *UserRepository) UpdateUser(ctx context.Context, user *models.User) error {
	query := `
		UPDATE users 
		SET name = $1, middle_name = $2, email = $3, phone = $4, role = $5, updated_at = NOW()
		WHERE id = $6 AND deleted_at IS NULL
	`
	_, err := repo.DB.Pool().Exec(ctx, query,
		user.Name,
		user.MiddleName,
		user.Email,
		user.Phone,
		user.Role,
		user.ID,
	)
	return err
}

func (repo *UserRepository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE users 
		SET deleted_at = NOW(), updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`
	_, err := repo.DB.Pool().Exec(ctx, query, id)
	return err
}

func NewUserRepository(db *database.Postgres) *UserRepository {
	return &UserRepository{
		DB: db,
	}
}

const (
	idByEmail = "SELECT id FROM users WHERE email = $1"
	idByPhone = "SELECT id FROM users WHERE phone = $1"
)

func (r *UserRepository) GetIdByField(ctx context.Context, field, value string) (uuid.UUID, error) {
	var query string

	switch field {
	case "email":
		query = idByEmail
	case "phone":
		query = idByPhone
	default:
		return uuid.Nil, fmt.Errorf("invalid field")
	}

	var id uuid.UUID
	err := r.DB.Pool().QueryRow(ctx, query, value).Scan(&id)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}
