package services

import (
	"citiaps/golang-backend-template/models"
	"citiaps/golang-backend-template/utils"
	"errors"
	"fmt"
)

// postgres
// se deberian hacer las validaciones respectivas

func CreateUserServicePostgres(user *models.UserPostgres) (*models.UserPostgres, error) {
	err := utils.ValidateUserPostgresObject(user)
	if err != nil {
		return nil, err
	}

	email := user.Email
	existUser, err := GetUserByEmailPostgres(email)
	if existUser != nil {
		return nil, errors.New("Ya existe usuario con el email.")
	}

	hashedPassword := utils.GeneratePassword(user.Password)

	user.Password = hashedPassword

	id, err := userRepoPostgres.InsertOnePostgres(user)

	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	user.ID = id

	return user, nil
}

func GetAllUsersServicePostgres() ([]*models.UserPostgres, error) {
	users, err := userRepoPostgres.GetAllPostgres()

	if err != nil {
		return nil, err
	}
	return users, nil
}

func GetUserByEmailPostgres(email string) (*models.UserPostgres, error) {
	condition := models.UserPostgres{Email: email}
	user, err := userRepoPostgres.FindOnePostgres(condition)

	if user == nil {
		return nil, errors.New("Error buscar usuario con email: " + email + ". No existe.")
	}

	if err != nil {
		return nil, err
	}
	return user, nil
}
