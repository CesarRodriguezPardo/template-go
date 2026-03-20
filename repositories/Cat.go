package repositories

import (
	"citiaps/golang-backend-template/models"

	"gorm.io/gorm"
)

// postgres
type CatRepositoryPostgres struct {
	DB *gorm.DB
}

// insertar un gatito
func (catRepoPostgres *CatRepositoryPostgres) InsertOnePostgres(cat *models.CatPostgres) (uint, error) {
	result := catRepoPostgres.DB.Create(cat)

	if result.Error != nil {
		return 0, result.Error
	}

	return cat.ID, nil
}

// buscar todos los gatitos
func (catRepoPostgres *CatRepositoryPostgres) FindAllPostgres() ([]*models.CatPostgres, error) {
	var cats []*models.CatPostgres
	result := catRepoPostgres.DB.Find(&cats)

	if result.Error != nil {
		return nil, result.Error
	}

	return cats, nil
}

// crear un cat repo
func NewCatRepositoryPostgres() *CatRepositoryPostgres {
	return &CatRepositoryPostgres{
		DB: DBPostgres.PostgresDB,
	}
}
