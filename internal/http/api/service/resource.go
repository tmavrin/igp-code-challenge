package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/tmavrin/igp-code-challenge/internal/store"
	"github.com/tmavrin/igp-code-challenge/internal/store/postgresdb"
)

type Resource struct {
	Fiber *fiber.App
	DB    store.Persistent
	Close func() error
}

func InitResource(ctx context.Context, config Config) (Resource, error) {
	var r Resource

	pool, closer, err := postgresdb.PgxPool(ctx, config.DatabaseURI)
	if err != nil {
		return r, fmt.Errorf("failed to initialize Postgres connection pool: %w", err)
	}

	r.DB = postgresdb.New(pool)

	r.Close = func() error {
		return errors.Join(
			closer(),
		)
	}

	return r, nil
}
