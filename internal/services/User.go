package services

import (
	"CesarRodriguezPardo/template-go/internal/models"
	"CesarRodriguezPardo/template-go/utils"
	"context"
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
)

func validateAndCapitaliceUser(user *models.User) error {
	if err := utils.ValidateMail(user.Email); err != nil {
		return fmt.Errorf("error validating mail: %w", err)
	}

	if err := utils.ValidatePhone(user.Email); err != nil {
		return fmt.Errorf("error validating phone number: %w", err)
	}

	if err := utils.ValidateString(user.Name); err != nil {
		return fmt.Errorf("error validating name: %w", err)
	}

	if err := utils.ValidateString(user.MiddleName); err != nil {
		return fmt.Errorf("error validating middle name: %w", err)
	}

	capitalizedName := utils.CapitalizateText(user.Name)
	capitalizedMiddleName := utils.CapitalizateText(user.MiddleName)

	user.Name = capitalizedName
	user.MiddleName = capitalizedMiddleName

	return nil
}

func findUserByEmail(ctx context.Context, email string) (uuid.UUID, error) {
	id, err := userRepo.FindUserByEmail(ctx, email)

	if err != nil {
		return uuid.Nil, fmt.Errorf("could not find user with email: %w", err)
	}

	return id, nil
}

func CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	if err := validateAndCapitaliceUser(user); err != nil {
		return nil, fmt.Errorf("error validating user: %w", err)
	}

	id, err := findUserByEmail(ctx, user.Email)

	if id {
		return nil, errors.New("user already exists with email")
	}

	return user, nil
}

func CreateUserServicePostgres(user *models.User) (*models.User, error) {
	err := utils.ValidateUserPostgresObject(user)
	if err != nil {
		return nil, err
	}

	email := user.Email
	existUser, err := GetUserByEmailPostgres(email)
	if existUser != nil {
		return nil, errors.New("Ya existe usuario con el email.")
	}

	hashedPassword, err := utils.GenerateHashedPassword(user.Password)

	if err != nil {
		return nil, fmt.Errorf("Could not hash password: %w", err)
	}

	user.Password = hashedPassword

	id, err := userRepoPostgres.InsertOnePostgres(user)

	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	user.ID = id

	return user, nil
}

func GetAllUsersServicePostgres() ([]*models.User, error) {
	users, err := userRepoPostgres.GetAllPostgres()

	if err != nil {
		return nil, err
	}
	return users, nil
}

func GetUserByEmailPostgres(email string) (*models.User, error) {
	condition := models.User{Email: email}
	user, err := userRepoPostgres.FindOnePostgres(condition)

	if user == nil {
		return nil, errors.New("Error buscar usuario con email: " + email + ". No existe.")
	}

	if err != nil {
		return nil, err
	}
	return user, nil
}
