package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/tmavrin/igp-code-challenge/internal/types"
)

type Tx interface {
	WithTx(ctx context.Context) (Persistent, error)
	CommitTx(ctx context.Context) error
	RollbackTx(ctx context.Context) error
}

type AccountsManager interface {
	AccountCreate(ctx context.Context, acc types.AuthCredentials) (types.Account, error)
	AccountGetBy(ctx context.Context, filter types.AccountFilter) (types.Account, error)
}

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -o ../fake/persistent_store.go --fake-name PersistentStoreProvider . Persistent
type Persistent interface {
	Tx
	AccountsManager
}

func IsErrNotFound(err error) bool {
	return errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows)
}
