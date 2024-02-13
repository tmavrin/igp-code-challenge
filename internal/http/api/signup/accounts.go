package signup

import (
	"github.com/gofiber/fiber/v2"
	"github.com/tmavrin/igp-code-challenge/internal/component/accounts"
)

type accountsRouter struct {
	component accounts.Provider
}

func NewAccountsRouter(component accounts.Provider) *accountsRouter {
	return &accountsRouter{component: component}
}

// Search godoc
//
//	@Summary		Search
//	@Description	Search accounts
//	@Tags			Account
//	@Accept			json
//	@Produce		json
//	@Failure		400
//	@Failure		500
//	@Router			/accounts/ [GET]
func (a *accountsRouter) Get(f *fiber.Ctx) error {

	f.SendString("GET ACCOUNTS")

	return nil
}
