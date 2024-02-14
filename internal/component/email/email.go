package email

import (
	"context"

	"github.com/gofiber/fiber/v2/log"
	"github.com/tmavrin/igp-code-challenge/internal/types"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -o ../../fake/component/email.go --fake-name EmailProvider . Provider
type Provider interface {
	SendWelcomeEmail(ctx context.Context, account types.Account) error
}

type component struct {
}

var _ Provider = (*component)(nil)

func New() *component {
	return &component{}
}

// This is a mock function that would send a "Welcome" email to user upon registering an account
func (c *component) SendWelcomeEmail(ctx context.Context, account types.Account) error {
	log.Debugf("sending welcome email to: %s", account.Email)

	// Can be sent with a client that would be initialized inside this servicer
	// c.emailClient.SendEmail(...)

	// Or it can be sent to another service trough Kafka or similar message queue
	// c.kafka.QueueMessage(...)
	return nil
}
