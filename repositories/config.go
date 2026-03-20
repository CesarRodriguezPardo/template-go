package repositories

import "citiaps/golang-backend-template/config/database"

const (
	CollectionCats  = "Cats"
	CollectionUsers = "Users"
)

var (
	DB         *database.MongoConnection
	DBPostgres *database.PostgresConnection
)

func InitConnections() {
	DB = database.NewMongoConnection()
	DBPostgres = database.NewPostgresConnection()
}
