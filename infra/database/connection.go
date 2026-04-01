package database

import (
	"CesarRodriguezPardo/template-go/config"
	"context"
	"fmt"
	"net/url"
	"sync"

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

func NewPG(ctx context.Context) (*Postgres, error) {
	uri := buildUri()

	pgOnce.Do(func() {
		db, err := pgxpool.New(ctx, uri)
		if err != nil {
			initErr = fmt.Errorf("pgxpool.New: %w", err)
			return
		}

		pgInstance = &Postgres{db}
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
