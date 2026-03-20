package models

import (
	"gorm.io/gorm"
)

// utilizando postgres
// gorm model trae algunas cositas interesantes para no repetir tanto codigo, este trae:
// ID uuint `gorm:"primaryKey"
// CreatedAt time.Time
// UpdatedAt time.Time
// DeletedAt gorm.DeteledAt

type CatPostgres struct {
	gorm.Model
	Name  string `json:"name" bson:"name,omitempty"`
	Age   int    `json:"age" bson:"age,omitempty"`
	Owner uint   `gorm:"foreignKey:UserRefer" json:"owner" bson:"string,omitempty"`
}
