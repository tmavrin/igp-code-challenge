package postgresdb

import (
	"context"
	"errors"

	"github.com/tmavrin/igp-code-challenge/internal/store"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Querier interface {
	Exec(ctx context.Context, sql string, arguments ...any) (commandTag pgconn.CommandTag, err error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
}

func New(db *pgxpool.Pool) *Queries {
	return &Queries{db: db}
}

type Queries struct {
	db Querier
}

func (q *Queries) WithTx(ctx context.Context) (store.Persistent, error) {
	db, ok := q.db.(*pgxpool.Pool)
	if !ok {
		return nil, errors.New("db not of type *pgxpool.Pool")
	}

	tx, err := db.Begin(ctx)
	if err != nil {
		return nil, err
	}

	return &Queries{
		db: tx,
	}, nil
}

func (q *Queries) CommitTx(ctx context.Context) error {
	return q.completeTx(ctx, false)
}

func (q *Queries) RollbackTx(ctx context.Context) error {
	return q.completeTx(ctx, true)
}

func (q *Queries) completeTx(ctx context.Context, rollback bool) error {
	tx, ok := q.db.(pgx.Tx)
	if !ok {
		return errors.New("db not of type pgx.Tx")
	}

	if rollback {
		return tx.Rollback(ctx)
	}

	return tx.Commit(ctx)
}

func PgxPool(ctx context.Context, uri string) (*pgxpool.Pool, func() error, error) {
	config, err := pgxpool.ParseConfig(uri)
	if err != nil {
		return nil, nil, err
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, nil, err
	}

	return pool,
		func() (err error) {
			pool.Close()
			return
		},
		pool.Ping(ctx)
}
