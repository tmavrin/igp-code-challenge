package accounts

import (
	"context"

	"github.com/tmavrin/igp-code-challenge/internal/store"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -o ../../../fake/component/accounts.go --fake-name AccountsProvider . Provider
type Provider interface {
	Get(ctx context.Context, query string) error
}

type component struct {
	persistent store.Persistent
}

var _ Provider = (*component)(nil)

func New(persistent store.Persistent) *component {
	return &component{persistent: persistent}
}

func (c *component) Get(ctx context.Context, query string) error {
	return nil
}
