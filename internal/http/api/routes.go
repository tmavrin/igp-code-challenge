package api

import (
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/tmavrin/igp-code-challenge/internal/component/auth"
	"github.com/tmavrin/igp-code-challenge/internal/component/email"
	"github.com/tmavrin/igp-code-challenge/internal/http/api/signup"
)

// sets up routes in separate file for better code readability
func (s *signupServer) routes() error {
	r := s.fiberApp

	emailComponent := email.New()
	authComponent := auth.New(
		s.Resource.DB,
		s.Config.Auth.JWTKey,
		s.Config.Auth.JWTDuration,
		emailComponent,
	)

	authRouter := signup.NewAuthRouter(authComponent)

	r.Post("/auth/register", authRouter.Register)
	r.Post("/auth/login", authRouter.Login)

	r.Use("/docs/*", basicauth.New(basicauth.Config{
		Users: map[string]string{
			s.Config.Docs.Username: s.Config.Docs.Password,
		},
	}))

	r.Static("/docs/", s.Config.Docs.Dir)

	return nil
}
