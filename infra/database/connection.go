package database

import (
	"CesarRodriguezPardo/template-go/config"
	"CesarRodriguezPardo/template-go/utils"
	"context"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/jackc/pgx/v5"
)

type Connection struct {
	Connect *pgx.Conn
}

func NewPostgresConnection() *Connection {
	ctx, cancel := context.WithTimeout(10 * time.Second)
	defer cancel()

	uri := buildUri()

	connCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	conn, err := pgx.Connect(connCtx, uri)

	if err != nil {
		utils.Fatal("Could not connect to Postgres.", errors.New("Could not connect to Postgres."))
	}
	defer conn.Close(ctx)

	connection := &Connection{}

	var ok error

	return connection
}

func buildUri() string {
	user := config.Cfg.Database.User
	pass := config.Cfg.Database.Pass
	host := config.Cfg.Database.Host
	port := config.Cfg.Database.Port
	dbName := config.Cfg.Database.Name

	uri := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(user, pass),
		Host:   fmt.Sprintf("%s:%s", host, port),
		Path:   dbName,
	}

	return uri.String()
}
