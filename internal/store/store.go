package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgx/v5"
)

type Tx interface {
	WithTx(ctx context.Context) (Persistent, error)
	CommitTx(ctx context.Context) error
	RollbackTx(ctx context.Context) error
}

// type NumberRateManager interface {
// 	ImportNumberRates(ctx context.Context, prefixes []string, values []float64) error
// 	DeleteNumberRates(ctx context.Context) error
// 	ArchiveNumberRates(ctx context.Context) error
// }

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -o ../fake/persistent_store.go --fake-name PersistentStoreProvider . Persistent
type Persistent interface {
	Tx
}

func IsErrNotFound(err error) bool {
	return errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows)
}
