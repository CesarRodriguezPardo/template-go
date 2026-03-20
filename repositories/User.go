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

// mongo
type UserRepository struct {
	Collection *mongo.Collection
}

// crear usuario
func (userRepo *UserRepository) InsertOne(user *models.User) (primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := userRepo.Collection.InsertOne(ctx, user)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return result.InsertedID.(primitive.ObjectID), nil
}

// obtener todos los usuarios
func (userRepo *UserRepository) GetAll() ([]*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.D{}
	cursor, err := userRepo.Collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var users []*models.User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, err
	}

	return users, nil
}

// obtener un usuario por filtro
func (userRepo *UserRepository) FindOne(filter bson.M) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user *models.User
	err := userRepo.Collection.FindOne(ctx, filter).Decode(&user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		Collection: DB.GetCollection(CollectionUsers),
	}
}

func (userRepo *UserRepository) CreateIndexes() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	indexModel := mongo.IndexModel{
		Keys: bson.M{
			"name": "text",
		},
		Options: options.Index().SetUnique(false),
	}

	_, err := userRepo.Collection.Indexes().CreateOne(ctx, indexModel)
	return err
}

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
