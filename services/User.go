package services

import (
	"citiaps/golang-backend-template/models"
	"citiaps/golang-backend-template/utils"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateUserService(user *models.User) (*models.User, error) {
	email := user.Email
	loweredEmail := utils.ToLowerEmail(email)
	user.Email = loweredEmail

	err := utils.ValidateUserObject(user)
	if err != nil {
		return nil, err
	}

	existUser, err := GetUserByEmail(loweredEmail)
	if existUser != nil {
		return nil, errors.New("Ya existe usuario con el email.")
	}

	time := time.Now()
	user.CreatedAt = time

	hashedPassword := utils.GeneratePassword(user.Password)
	user.Password = hashedPassword

	id, err := userRepo.InsertOne(user)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	user.ID = id

	return user, nil
}

func GetAllUsersService() ([]*models.User, error) {
	users, err := userRepo.GetAll()

	if err != nil {
		return nil, err
	}
	return users, nil
}

func GetUserByEmail(email string) (*models.User, error) {
	filter := bson.M{"email": email}
	user, err := userRepo.FindOne(filter)

	if user == nil {
		return nil, errors.New("Error buscar user con email: " + email + ". No existe.")
	}

	if err != nil {
		return nil, err
	}
	return user, nil
}

func GetUserById(id primitive.ObjectID) (*models.User, error) {
	filter := bson.M{"_id": id}
	user, err := userRepo.FindOne(filter)

	if user == nil {
		return nil, errors.New("Error al buscar usuario con ID. No existe")
	}

	if err != nil {
		return nil, err
	}
	return user, nil
}

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
