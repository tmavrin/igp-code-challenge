//go:build integration

package postgresdb

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/require"
	"github.com/tmavrin/igp-code-challenge/internal/types"
)

func TestCreateAccount(t *testing.T) {
	defer truncate()

	var (
		databaseManager = New(testDB)
	)

	type args struct {
		credentials types.AuthCredentials
	}
	tests := []struct {
		name              string
		args              args
		expectedOutput    types.Account
		expectedErrorCode string
	}{
		{
			name: "it should create account",
			args: args{
				credentials: types.AuthCredentials{
					Email:    "test@mail.com",
					Password: "some-hash",
				},
			},
			expectedOutput: types.Account{
				Email: "test@mail.com",
			},
		},
		{
			name: "it should return unique constraint violation on duplicate email",
			args: args{
				credentials: types.AuthCredentials{
					Email:    "test@mail.com",
					Password: "some-hash",
				},
			},
			expectedErrorCode: pgerrcode.UniqueViolation,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			createdAccount, err := databaseManager.AccountCreate(context.Background(), tt.args.credentials)

			if err != nil {
				var pgErr *pgconn.PgError
				require.True(t, errors.As(err, &pgErr) && tt.expectedErrorCode == pgErr.Code)
			}

			if tt.expectedErrorCode == "" {
				require.Equal(t, tt.expectedOutput.Email, createdAccount.Email)
				require.NotEmpty(t, createdAccount.ID)
				require.NotEmpty(t, createdAccount.Created)
				require.NotEmpty(t, createdAccount.Updated)
			}
		})
	}
}

func TestGetAccountBy(t *testing.T) {
	defer truncate()

	var (
		databaseManager = New(testDB)
		pass            = "some-hash"
		invalidEmail    = "invalid@mail.com"
	)

	acc1 := types.Account{
		ID:       uuid.New(),
		Email:    "test1@mail.com",
		Password: &pass,
	}
	acc2 := types.Account{
		ID:       uuid.New(),
		Email:    "test2@mail.com",
		Password: &pass,
	}

	type args struct {
		filter types.AccountFilter
	}
	tests := []struct {
		name           string
		args           args
		expectedOutput types.Account
		expectedError  error
		prepare        func() error
	}{
		{
			name: "it should get account by id",
			prepare: func() error {
				_, err := databaseManager.db.Exec(context.Background(),
					"INSERT INTO accounts (id, email, password_hash) VALUES ($1, $2, $3)",
					acc1.ID, acc1.Email, acc1.Password,
				)
				return err
			},
			args: args{
				filter: types.AccountFilter{
					ByID: &acc1.ID,
				},
			},
			expectedOutput: acc1,
		},
		{
			name: "it should get other account by email",
			prepare: func() error {
				_, err := databaseManager.db.Exec(context.Background(),
					"INSERT INTO accounts (id, email, password_hash) VALUES ($1, $2, $3)",
					acc2.ID, acc2.Email, acc2.Password,
				)
				return err
			},
			args: args{
				filter: types.AccountFilter{
					ByEmail: &acc2.Email,
				},
			},
			expectedOutput: acc2,
		},
		{
			name: "it should return error on not found",
			prepare: func() error {
				return nil
			},
			args: args{
				filter: types.AccountFilter{
					ByEmail: &invalidEmail,
				},
			},
			expectedError: pgx.ErrNoRows,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			err := tt.prepare()
			require.NoError(t, err)

			account, err := databaseManager.AccountGetBy(context.Background(), tt.args.filter)

			require.ErrorIs(t, err, tt.expectedError)

			if tt.expectedError == nil {
				require.Equal(t, tt.expectedOutput.Email, account.Email)
			}
		})
	}
}
