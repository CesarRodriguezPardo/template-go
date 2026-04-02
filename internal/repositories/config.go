package repositories

import (
	"CesarRodriguezPardo/template-go/infra/database"
	"context"
	"fmt"
)

const (
	CollectionCats  = "Cats"
	CollectionUsers = "Users"
)

var (
	DB *database.Postgres
)

func InitConnections() error {
	var err error

	DB, err = database.Connect(context.Background())

	if err != nil {
		return fmt.Errorf("could not connect to PG: %w", err)
	}

	return nil
}
