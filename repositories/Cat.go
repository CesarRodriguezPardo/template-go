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

// en este caso el ejemplo de repositorio va a ser usando mongo directamente.
type CatRepository struct {
	Collection *mongo.Collection
}

// insertar un gatito
func (catRepo *CatRepository) InsertOne(cat *models.Cat) (primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := catRepo.Collection.InsertOne(ctx, cat)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return result.InsertedID.(primitive.ObjectID), nil
}

// buscar un gatito
func (catRepo *CatRepository) FindOne(filter bson.M) (*models.Cat, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var cat *models.Cat

	err := catRepo.Collection.FindOne(ctx, filter).Decode(&cat)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return cat, nil
}

// buscar todos los gatitos
func (catRepo *CatRepository) FindAll() ([]*models.Cat, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{}
	cursor, err := catRepo.Collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var cats []*models.Cat
	if err := cursor.All(ctx, &cats); err != nil {
		return nil, err
	}

	return cats, nil
}

// buscar todos los gatitos por filtro
func (catRepo *CatRepository) FindAllFiltered(filter bson.D) ([]*models.Cat, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := catRepo.Collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var cats []*models.Cat
	if err := cursor.All(ctx, &cats); err != nil {
		return nil, err
	}

	return cats, nil
}

func NewCatRepository() *CatRepository {
	return &CatRepository{
		Collection: DB.GetCollection(CollectionCats),
	}
}

func (catRepo *CatRepository) CreateIndexes() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	indexModel := mongo.IndexModel{
		Keys: bson.M{
			"name": "text",
		},
		Options: options.Index().SetUnique(false),
	}

	_, err := catRepo.Collection.Indexes().CreateOne(ctx, indexModel)
	return err
}

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
