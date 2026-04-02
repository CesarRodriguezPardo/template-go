package database

import (
	"CesarRodriguezPardo/template-go/config"
	"context"
	"fmt"
	"net/url"
	"sync"

	logger "CesarRodriguezPardo/template-go/infra/logger"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	db *pgxpool.Pool
}

var (
	pgInstance *Postgres
	pgOnce     sync.Once
)

var initErr error

func Connect(ctx context.Context) (*Postgres, error) {
	uri := buildUri()

	pgOnce.Do(func() {
		db, err := pgxpool.New(ctx, uri)
		if err != nil {
			initErr = fmt.Errorf("pgxpool.New: %w", err)
			return
		}

		if err = db.Ping(ctx); err != nil {
			initErr = fmt.Errorf("db.Ping: %w", err)
			return
		}

		pgInstance = &Postgres{db}

		connectData := getConnectionData()
		logger.Info("Database - " + connectData)
	})

	return pgInstance, initErr
}

func (pg *Postgres) Pool() *pgxpool.Pool {
	return pg.db
}

func (pg *Postgres) Ping(ctx context.Context) error {
	return pg.db.Ping(ctx)
}

func (pg *Postgres) Close() {
	pg.db.Close()
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

func getConnectionData() string {
	host := config.Cfg.Database.Host
	port := config.Cfg.Database.Port
	data := fmt.Sprintf("host: %s, port: %s", host, port)
	return data
}
