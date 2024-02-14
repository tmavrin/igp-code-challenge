package types

import (
	"time"

	"github.com/google/uuid"
)

type Account struct {
	ID       uuid.UUID `json:"id"`
	Email    string    `json:"email"`
	Password *string   `json:"-"`
	Created  time.Time `json:"created"`
	Updated  time.Time `json:"updated"`
}

type AccountFilter struct {
	ByID    *uuid.UUID
	ByEmail *string
}
