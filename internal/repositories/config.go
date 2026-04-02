package repositories

import (
	"CesarRodriguezPardo/template-go/infra/database"
	"context"
	"fmt"
)

const (
	CollectionUsers = "Users"
)

var (
	DB *database.Postgres
)

func InitConnections(ctx context.Context) (*database.Postgres, error) {
	db, err := database.Connect(ctx)

	if err != nil {
		return nil, fmt.Errorf("could not connect to PG: %w", err)
	}

	return db, nil
}
