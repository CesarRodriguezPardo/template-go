package repositories

import "citiaps/golang-backend-template/config/database"

const (
	CollectionCats  = "Cats"
	CollectionUsers = "Users"
)

var (
	DBPostgres *database.PostgresConnection
)

func InitConnections() {
	DBPostgres = database.NewPostgresConnection()
}
