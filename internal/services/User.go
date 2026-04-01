package services

import (
	"CesarRodriguezPardo/template-go/internal/models"
	"CesarRodriguezPardo/template-go/utils"
	"context"
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
)

func validateUserParams(user *models.User) error {
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

	return nil
}

func capitaliceUserParams(user *models.User) error {
	capitalizedName := utils.CapitalizateText(user.Name)
	capitalizedMiddleName := utils.CapitalizateText(user.MiddleName)

	user.Name = capitalizedName
	user.MiddleName = capitalizedMiddleName

	return nil
}

func findUserByEmail(ctx context.Context, email string) (uuid.UUID, error) {
	id, err := userRepo.GetIdByEmail(ctx, email)

	if err != nil {
		return uuid.Nil, fmt.Errorf("could not find user with email: %w", err)
	}

	return id, nil
}

func CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	if err := validateUserParams(user); err != nil {
		return nil, fmt.Errorf("error validating user: %w", err)
	}
	capitaliceUserParams(user)

	toFindId, err := findUserByEmail(ctx, user.Email)
	if toFindId != uuid.Nil {
		return nil, errors.New("user already exists with email")
	}
	if err != nil {
		return nil, fmt.Errorf("could not find user with email: %w", err)
	}

	hashedPass, err := utils.GenerateHashedPassword(user.Password)
	if err != nil {
		return nil, fmt.Errorf("could not hash password: %w", err)
	}

	user.Password = hashedPass

	id, err := userRepo.CreateUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("could not create user: %w", err)
	}

	user.ID = id
	return user, nil
}

// buscar todos los usuarios
