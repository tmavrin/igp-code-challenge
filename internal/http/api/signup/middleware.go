package signup

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/tmavrin/igp-code-challenge/internal/component/auth"
)

func Auth(component auth.Provider) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		token := c.GetReqHeaders()["Authorization"][0]
		if token == "" {
			return fiber.ErrUnauthorized
		}
		parts := strings.Split(token, " ")
		if len(parts) != 2 {
			return fiber.ErrUnauthorized
		}

		account, err := component.Auth(c.Context(), parts[1])
		if err != nil {
			return fiber.ErrUnauthorized
		}

		c.Locals("account", account)

		return c.Next()
	}

}
