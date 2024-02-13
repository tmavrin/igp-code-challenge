package api

import (
	"github.com/tmavrin/igp-code-challenge/internal/component/accounts"
	"github.com/tmavrin/igp-code-challenge/internal/http/api/signup"
)

// sets up routes in separate file for better code readability
func (s *signupServer) routes() error {
	r := s.fiberApp

	accountsComponent := accounts.New(s.Resource.DB)

	accountsRouter := signup.NewAccountsRouter(accountsComponent)

	r.Get("/accs", accountsRouter.Get)

	return nil
}
