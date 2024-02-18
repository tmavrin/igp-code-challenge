package auth

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	emailComponent "github.com/tmavrin/igp-code-challenge/internal/component/email"
	"github.com/tmavrin/igp-code-challenge/internal/fake"
	"github.com/tmavrin/igp-code-challenge/internal/types"
	"golang.org/x/crypto/bcrypt"
)

func TestLogin(t *testing.T) {
	var (
		jwtKey      = "test_key"
		jwtDuration = time.Duration(time.Hour)
		email       = "test@mail.com"
		password    = "password"
	)

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	require.NoError(t, err)

	pass := string(bytes)

	acc := types.Account{
		ID:       uuid.New(),
		Email:    email,
		Password: &pass,
	}

	type fields struct {
		persistentStoreProvider *fake.PersistentStoreProvider
	}
	type args struct {
		credentials types.AuthCredentials
	}
	tests := []struct {
		name            string
		fields          fields
		args            args
		expectedError   error
		expectedAccount types.Account
	}{
		{
			name: "it should return user and token on success",
			fields: fields{
				persistentStoreProvider: &fake.PersistentStoreProvider{
					AccountGetByStub: func(ctx context.Context, af types.AccountFilter) (types.Account, error) {
						return acc, nil
					},
				},
			},
			args: args{
				credentials: types.AuthCredentials{Email: email, Password: password},
			},
			expectedError:   nil,
			expectedAccount: acc,
		},
		{
			name: "it should return user error on invalid password",
			fields: fields{
				persistentStoreProvider: &fake.PersistentStoreProvider{
					AccountGetByStub: func(ctx context.Context, af types.AccountFilter) (types.Account, error) {
						return acc, nil
					},
				},
			},
			args: args{
				credentials: types.AuthCredentials{Email: email, Password: "bad pass"},
			},
			expectedError:   ErrorInvalidCredentials,
			expectedAccount: types.Account{},
		},
		{
			name: "it should return user error on user not found",
			fields: fields{
				persistentStoreProvider: &fake.PersistentStoreProvider{
					AccountGetByStub: func(ctx context.Context, af types.AccountFilter) (types.Account, error) {
						return types.Account{}, pgx.ErrNoRows
					},
				},
			},
			args: args{
				credentials: types.AuthCredentials{Email: email, Password: pass},
			},
			expectedError:   ErrorUserNotFound,
			expectedAccount: types.Account{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New(tt.fields.persistentStoreProvider, []byte(jwtKey), jwtDuration, nil)
			account, token, err := c.Login(context.Background(), tt.args.credentials)
			require.ErrorIs(t, err, tt.expectedError)

			if tt.expectedError != nil {
				require.Empty(t, token)
			} else {
				require.NotEmpty(t, token)
			}

			require.EqualValues(t, tt.expectedAccount, account)
		})
	}
}

func TestCreate(t *testing.T) {
	var (
		jwtKey      = "test_key"
		jwtDuration = time.Duration(time.Hour)
		email       = "test@mail.com"
		password    = "password"
	)

	type fields struct {
		persistentStoreProvider *fake.PersistentStoreProvider
	}
	type args struct {
		credentials types.AuthCredentials
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		expectedError error
	}{
		{
			name: "it should succeed",
			fields: fields{
				persistentStoreProvider: &fake.PersistentStoreProvider{
					AccountCreateStub: func(ctx context.Context, ac types.AuthCredentials) (types.Account, error) {
						return types.Account{}, nil
					},
				},
			},
			args: args{
				credentials: types.AuthCredentials{Email: email, Password: password},
			},
			expectedError: nil,
		},
		{
			name: "it should return error on duplicate",
			fields: fields{
				persistentStoreProvider: &fake.PersistentStoreProvider{
					AccountCreateStub: func(ctx context.Context, ac types.AuthCredentials) (types.Account, error) {
						return types.Account{}, &pgconn.PgError{Code: pgerrcode.UniqueViolation}
					},
				},
			},
			args: args{
				credentials: types.AuthCredentials{Email: email, Password: password},
			},
			expectedError: ErrorUserExists,
		},
		{
			name: "it should return error on provider error",
			fields: fields{
				persistentStoreProvider: &fake.PersistentStoreProvider{
					AccountCreateStub: func(ctx context.Context, ac types.AuthCredentials) (types.Account, error) {
						return types.Account{}, assert.AnError
					},
				},
			},
			args: args{
				credentials: types.AuthCredentials{Email: email, Password: password},
			},
			expectedError: assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New(tt.fields.persistentStoreProvider, []byte(jwtKey), jwtDuration, emailComponent.New())
			err := c.Create(context.Background(), tt.args.credentials)
			require.ErrorIs(t, err, tt.expectedError)
		})
	}
}
