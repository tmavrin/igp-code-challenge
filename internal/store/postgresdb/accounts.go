package postgresdb

import (
	"context"
	"fmt"
	"strings"

	"github.com/tmavrin/igp-code-challenge/internal/types"
)

func (q *Queries) AccountCreate(ctx context.Context, credentials types.AuthCredentials) (types.Account, error) {
	acc := types.Account{Email: credentials.Email}
	query := "INSERT INTO accounts (email, password_hash) VALUES ($1, $2) RETURNING id, created, updated"
	err := q.db.QueryRow(ctx, query, credentials.Email, credentials.Password).Scan(
		&acc.ID,
		&acc.Created,
		&acc.Updated,
	)

	return acc, err
}

func (q *Queries) AccountGetBy(ctx context.Context, filter types.AccountFilter) (types.Account, error) {
	var (
		whereClause []string
		args        []any
	)

	if filter.ByID != nil {
		args = append(args, filter.ByID)
		whereClause = append(whereClause, fmt.Sprintf("id = $%d", len(args)))
	}

	if filter.ByEmail != nil {
		args = append(args, filter.ByEmail)
		whereClause = append(whereClause, fmt.Sprintf("email = $%d", len(args)))
	}

	if len(args) == 0 {
		return types.Account{}, ErrorNoFiltersProvided
	}

	query := fmt.Sprintf(`
		SELECT
			id,
			email,
			password_hash,
			created,
			updated
		FROM accounts
		WHERE
			%s`,
		strings.Join(whereClause, " AND "),
	)

	var acc types.Account
	err := q.db.QueryRow(ctx, query, args...).Scan(
		&acc.ID,
		&acc.Email,
		&acc.Password,
		&acc.Created,
		&acc.Updated,
	)

	return acc, err
}
