package services

import (
	"CesarRodriguezPardo/template-go/internal/models"
	"CesarRodriguezPardo/template-go/utils"
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
	uuid "github.com/satori/go.uuid"
)

func validateUserParams(user *models.User) error {
	if err := utils.ValidateMail(user.Email); err != nil {
		return fmt.Errorf("error validating mail: %w", err)
	}

	if err := utils.ValidatePhone(user.Phone); err != nil {
		return fmt.Errorf("error validating phone number: %w", err)
	}

	if err := utils.ValidateString(user.Name); err != nil {
		return fmt.Errorf("error validating name: %w", err)
	}

	if err := utils.ValidateString(user.MiddleName); err != nil {
		return fmt.Errorf("error validating middle name: %w", err)
	}

	return nil
}

func capitaliceUserParams(user *models.User) {
	capitalizedName := utils.CapitalizateText(user.Name)
	capitalizedMiddleName := utils.CapitalizateText(user.MiddleName)

	user.Name = capitalizedName
	user.MiddleName = capitalizedMiddleName
}

func findUserByField(ctx context.Context, field, value string) (uuid.UUID, error) {
	id, err := userRepo.GetIdByField(ctx, field, value)

	if err != nil {
		return uuid.Nil, fmt.Errorf("could not find user with field %s: %w", field, err)
	}

	return id, nil
}

/*
func findUserByEmail(ctx context.Context, email string) (uuid.UUID, error) {
	id, err := userRepo.GetIdByEmail(ctx, email)

	if err != nil {
		return uuid.Nil, fmt.Errorf("could not find user with email: %w", err)
	}

	return id, nil
}
*/

func CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	if err := validateUserParams(user); err != nil {
		return nil, fmt.Errorf("error validating user: %w", err)
	}
	capitaliceUserParams(user)

	// Forzar rol por defecto — nunca confiar en el cliente
	user.Role = string(models.WORKER)

	hashedPass, err := utils.GenerateHash(user.Password)
	if err != nil {
		return nil, fmt.Errorf("could not hash password: %w", err)
	}
	user.Password = hashedPass

	id, err := userRepo.CreateUser(ctx, user)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			if strings.Contains(pgErr.ConstraintName, "email") {
				return nil, errors.New("user already exists with email")
			}
			if strings.Contains(pgErr.ConstraintName, "phone") {
				return nil, errors.New("user already exists with phone number")
			}
		}
		return nil, errors.New("could not create user")
	}

	user.ID = id
	user.Password = ""
	return user, nil
}

func GetAllUsers(ctx context.Context) ([]*models.User, error) {
	users, err := userRepo.GetAllUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not get users: %w", err)
	}

	return users, nil
}

// buscar todos los usuarios
