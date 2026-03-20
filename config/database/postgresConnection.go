package database

import (
	"citiaps/golang-backend-template/config"
	"citiaps/golang-backend-template/models"
	"citiaps/golang-backend-template/utils"
	"context"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresConnection struct {
	PostgresDB *gorm.DB
}

func NewPostgresConnection() *PostgresConnection {
	ctxTimeout := 10 * time.Second

	connection := &PostgresConnection{}

	var ok error

	postgresUri := config.Cfg.Database.Postgres.URI

	connection.PostgresDB, ok = gorm.Open(postgres.Open(postgresUri))
	if ok != nil {
		utils.Fatal("%w", ok)
		return nil
	}

	pgdb, err := connection.PostgresDB.DB()
	if err != nil {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()
	err = pgdb.PingContext(ctx)
	if err != nil {
		return nil
	}

	utils.Info("POSTGRES: Conexion exitosa.")

	// creacion de tablas
	err = connection.PostgresDB.AutoMigrate(&models.CatPostgres{}, &models.UserPostgres{})
	if err != nil {
		utils.Fatal("POSTGRES: %v", err)
	}

	utils.Info("POSTGRES: Creacion de tablas realizada.")
	return connection
}

/*

REFS:
	- https://gorm.io/docs/gorm_config.html
	- https://go.dev/wiki/SQLInterface
	- https://pkg.go.dev/database/sql#DB.PingContext
*/
