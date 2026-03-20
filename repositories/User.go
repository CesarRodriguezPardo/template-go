package repositories

import (
	"citiaps/golang-backend-template/models"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/gorm"
)

// postgres
type UserRepositoryPostgres struct {
	DB *gorm.DB
}

// crear usuario
func (userRepoPostgres *UserRepositoryPostgres) InsertOnePostgres(user *models.UserPostgres) (uint, error) {
	result := userRepoPostgres.DB.Create(user)

	if result.Error != nil {
		return 0, result.Error
	}

	return user.ID, nil
}

// obtener todos los usuarios
func (userRepoPostgres *UserRepositoryPostgres) GetAllPostgres() ([]*models.UserPostgres, error) {
	var allUsers []*models.UserPostgres

	result := userRepoPostgres.DB.Find(&allUsers)

	if result.Error != nil {
		return nil, result.Error
	}

	return allUsers, nil
}

func (UserRepoPostgres *UserRepositoryPostgres) FindOnePostgres(filter models.UserPostgres) (*models.UserPostgres, error) {
	var user *models.UserPostgres

	result := UserRepoPostgres.DB.First(&user, filter)

	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func NewUserRepositoryPostgres() *UserRepositoryPostgres {
	return &UserRepositoryPostgres{
		DB: DBPostgres.PostgresDB,
	}
}
