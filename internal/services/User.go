package services

import (
	"CesarRodriguezPardo/template-go/internal/models"
	"CesarRodriguezPardo/template-go/utils"
	"errors"
	"fmt"
)

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
