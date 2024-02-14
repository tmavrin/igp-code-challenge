package types

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AuthClaims struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
	jwt.RegisteredClaims
}

type AuthCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
