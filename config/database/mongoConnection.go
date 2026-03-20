package database

import (
	"citiaps/golang-backend-template/config"
	"citiaps/golang-backend-template/utils"
	"context"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoConnection struct {
	MongoClient *mongo.Client
}

func NewMongoConnection() *MongoConnection {

	ctxTimeout := 10 * time.Second

	connection := &MongoConnection{}

	var ok error

	mongoUri := config.Cfg.Database.Mongo.URI

	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()

	clientOptions := options.Client().ApplyURI(mongoUri)
	connection.MongoClient, ok = mongo.Connect(ctx, clientOptions)

	if ok != nil {
		utils.Fatal("error conectando a MongoDB: %w", ok)
	}

	err := connection.MongoClient.Ping(ctx, readpref.Primary())
	if err != nil {
		_ = connection.MongoClient.Disconnect(context.Background())
		utils.Fatal("error haciendo ping a MongoDB:", err)
	}

	utils.Info("MONGO: Conexion exitosa.")

	return connection
}

func (d *MongoConnection) GetCollection(collection string) *mongo.Collection {
	return d.MongoClient.Database(os.Getenv("DB_NAME_MONGO")).Collection(collection)
}

/*
REFS:
	- https://www.mongodb.com/docs/drivers/go/current/fundamentals/connections/connection-guide/
	- https://www.mongodb.com/docs/drivers/go/current/fundamentals/connections/connection-options/#std-label-golang-connection-options
*/
