package services

import (
	"citiaps/golang-backend-template/models"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)


// postgres

func CreateCatServicePostgres(cat *models.CatPostgres) (*models.CatPostgres, error) {
	if cat.Name == "" {
		return nil, errors.New("El nombre del gatito/a no puede ser vacio.")
	}
	if cat.Age <= 0 {
		return nil, errors.New("La edad del gatito/a no puede ser cero o negativo.")
	}

	// se asume que no hay gatitos en adopcion, es decir, sin dueño.
	if cat.Owner == 0 {
		return nil, errors.New("No existe dueño del gatito/a.")
	}

	id, err := catRepoPostgres.InsertOnePostgres(cat)

	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	cat.ID = id

	return cat, nil
}

func FindAllCatServicePostgres() ([]*models.CatPostgres, error) {
	cats, err := catRepoPostgres.FindAllPostgres()

	if err != nil {
		return nil, err
	}

	return cats, nil
}
