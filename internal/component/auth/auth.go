package auth

import (
	"context"
	"errors"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/tmavrin/igp-code-challenge/internal/component/email"
	"github.com/tmavrin/igp-code-challenge/internal/store"
	"github.com/tmavrin/igp-code-challenge/internal/types"
	"golang.org/x/crypto/bcrypt"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -o ../../fake/component/auth.go --fake-name AuthProvider . Provider
type Provider interface {
	Login(ctx context.Context, credentials types.AuthCredentials) (types.Account, string, error)
	Auth(ctx context.Context, token string) (types.Account, error)
	Create(ctx context.Context, account types.AuthCredentials) error
}

type component struct {
	persistent  store.Persistent
	jwtKey      []byte
	jwtDuration time.Duration
	email       email.Provider
}

var _ Provider = (*component)(nil)

func New(
	persistent store.Persistent,
	jwtKey []byte,
	jwtDuration time.Duration,
	email email.Provider,
) *component {
	return &component{
		persistent:  persistent,
		jwtKey:      jwtKey,
		jwtDuration: jwtDuration,
		email:       email,
	}
}

var (
	ErrorUserExists         = errors.New("user exists")
	ErrorUserNotFound       = errors.New("account not found")
	ErrorInvalidCredentials = errors.New("invalid credentials")
)

func (c *component) Login(ctx context.Context, credentials types.AuthCredentials) (types.Account, string, error) {
	account, err := c.persistent.AccountGetBy(ctx, types.AccountFilter{ByEmail: &credentials.Email})
	if errors.Is(err, pgx.ErrNoRows) {
		return types.Account{}, "", ErrorUserNotFound
	}
	if err != nil {
		return types.Account{}, "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(*account.Password), []byte(credentials.Password))
	if err != nil {
		return types.Account{}, "", ErrorInvalidCredentials
	}

	authClaims := types.AuthClaims{
		ID:    account.ID,
		Email: account.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.NewString(),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(c.jwtDuration)),
		},
	}
	tokenString, err := jwt.NewWithClaims(jwt.SigningMethodHS256, authClaims).SignedString(c.jwtKey)

	return account, tokenString, err
}

// To be used inside auth middleware that checks JWT token
func (c *component) Auth(ctx context.Context, token string) (types.Account, error) {
	var authClaims types.AuthClaims
	_, err := jwt.ParseWithClaims(token, &authClaims, func(token *jwt.Token) (interface{}, error) {
		return c.jwtKey, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}))
	if err != nil {
		return types.Account{}, err
	}

	return types.Account{
		ID:    authClaims.ID,
		Email: authClaims.Email,
	}, nil
}

func (c *component) Create(ctx context.Context, credentials types.AuthCredentials) error {

	bytes, err := bcrypt.GenerateFromPassword([]byte(credentials.Password), 12)
	if err != nil {
		return err
	}

	credentials.Password = string(bytes)

	account, err := c.persistent.AccountCreate(ctx, credentials)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgerrcode.UniqueViolation == pgErr.Code {
		return ErrorUserExists
	}
	if err != nil {
		return err
	}

	go func() {
		err := c.email.SendWelcomeEmail(ctx, account)
		if err != nil {
			log.Errorf("failed to send a welcome mail: %s", err)
		}
	}()

	return nil
}
