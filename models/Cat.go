package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gorm.io/gorm"
)

// utilizando mongo

type Cat struct {
	ID    primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Name  string             `json:"name" bson:"name,omitempty"`
	Age   int                `json:"age" bson:"age,omitempty"`
	Owner primitive.ObjectID `json:"owner" bson:"string,omitempty"`
}

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
