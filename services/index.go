package services

import (
	"citiaps/golang-backend-template/repositories"
	"errors"
)

var (
	// postgres
	catRepoPostgres  *repositories.CatRepositoryPostgres
	userRepoPostgres *repositories.UserRepositoryPostgres
)

func InitRepositories() {
	catRepoPostgres = repositories.NewCatRepositoryPostgres()
	userRepoPostgres = repositories.NewUserRepositoryPostgres()
}

func InitIndexes() error {
	err := catRepo.CreateIndexes()
	if err != nil {
		return errors.New("error al inicializar los indices de los gatitos: " + err.Error())
	}

	err = userRepo.CreateIndexes()
	if err != nil {
		return errors.New("error al inicializar los índices de los usuarios: " + err.Error())
	}
	return nil
}
