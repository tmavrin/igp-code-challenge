package service

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tmavrin/igp-code-challenge/internal/component/auth"
	"github.com/tmavrin/igp-code-challenge/internal/fake/component"
	"github.com/tmavrin/igp-code-challenge/internal/types"
)

func TestRegister(t *testing.T) {
	type fields struct {
		accountsProvider *component.AuthProvider
	}
	tests := []struct {
		name           string
		fields         fields
		body           string
		expectedCode   int
		expectedOutput string
	}{
		{
			name: "it should create account",
			fields: fields{
				accountsProvider: &component.AuthProvider{
					CreateStub: func(ctx context.Context, ac types.AuthCredentials) error {
						return nil
					},
				},
			},
			body:           `{ "email": "test1@email.com","password":"password"}`,
			expectedCode:   http.StatusOK,
			expectedOutput: "",
		},
		{
			name: "it should return bad request on bad email",
			fields: fields{
				accountsProvider: &component.AuthProvider{
					CreateStub: func(ctx context.Context, ac types.AuthCredentials) error {
						return nil
					},
				},
			},
			body:           `{ "email": "invalid-email","password":"password"}`,
			expectedCode:   http.StatusBadRequest,
			expectedOutput: "Bad Request",
		},
		{
			name: "it should return conflict on duplicate account",
			fields: fields{
				accountsProvider: &component.AuthProvider{
					CreateStub: func(ctx context.Context, ac types.AuthCredentials) error {
						return auth.ErrorUserExists
					},
				},
			},
			body:           `{ "email": "test1@email.com","password":"password"}`,
			expectedCode:   http.StatusConflict,
			expectedOutput: "Conflict",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := NewAuthRouter(tt.fields.accountsProvider)
			app := fiber.New()
			app.Post("/auth/register", router.Register)

			req := httptest.NewRequest("POST", "/auth/register", bytes.NewReader([]byte(tt.body)))
			req.Header.Add("Content-Type", "application/json")

			resp, err := app.Test(req)
			require.NoError(t, err)
			respBody, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			require.Equal(t, tt.expectedCode, resp.StatusCode)
			require.Equal(t, tt.expectedOutput, string(respBody))
		})
	}
}

func TestLogin(t *testing.T) {
	type fields struct {
		accountsProvider *component.AuthProvider
	}
	tests := []struct {
		name           string
		fields         fields
		body           string
		expectedCode   int
		expectedOutput string
	}{
		{
			name: "it should login",
			fields: fields{
				accountsProvider: &component.AuthProvider{
					LoginStub: func(ctx context.Context, ac types.AuthCredentials) (types.Account, string, error) {
						return types.Account{
							ID:    uuid.New(),
							Email: "test@mail.com",
						}, "token", nil
					},
				},
			},
			body:           `{ "email": "test@email.com","password":"password"}`,
			expectedCode:   http.StatusOK,
			expectedOutput: `{"id":".*","email":"test@mail.com","created":".*","updated":".*","token":"token"}`,
		},
		{
			name: "it should return unauthorized on bad login",
			fields: fields{
				accountsProvider: &component.AuthProvider{
					LoginStub: func(ctx context.Context, ac types.AuthCredentials) (types.Account, string, error) {
						return types.Account{}, "", auth.ErrorInvalidCredentials
					},
				},
			},
			body:           `{ "email": "unknown@email.com","password":"password"}`,
			expectedCode:   http.StatusUnauthorized,
			expectedOutput: "Unauthorized",
		},
		{
			name: "it should return 500 unknown component error",
			fields: fields{
				accountsProvider: &component.AuthProvider{
					LoginStub: func(ctx context.Context, ac types.AuthCredentials) (types.Account, string, error) {
						return types.Account{}, "", assert.AnError
					},
				},
			},
			body:           `{ "email": "unknown@email.com","password":"password"}`,
			expectedCode:   http.StatusInternalServerError,
			expectedOutput: "assert.AnError general error for testing",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := NewAuthRouter(tt.fields.accountsProvider)
			app := fiber.New()
			app.Post("/auth/login", router.Login)

			req := httptest.NewRequest("POST", "/auth/login", bytes.NewReader([]byte(tt.body)))
			req.Header.Add("Content-Type", "application/json")

			resp, err := app.Test(req)
			require.NoError(t, err)
			respBody, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			require.Equal(t, tt.expectedCode, resp.StatusCode)
			require.Regexp(t, regexp.MustCompile(tt.expectedOutput), string(respBody))
		})
	}
}
