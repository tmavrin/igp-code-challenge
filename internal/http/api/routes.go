package api

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/tmavrin/igp-code-challenge/internal/component/auth"
	"github.com/tmavrin/igp-code-challenge/internal/component/email"
	"github.com/tmavrin/igp-code-challenge/internal/component/notifications"
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
	notificationsComponent := notifications.New()

	authRouter := signup.NewAuthRouter(authComponent)
	notificationsRouter := signup.NewNotificationsRouter(notificationsComponent)

	r.Post("/auth/register", authRouter.Register)
	r.Post("/auth/login", authRouter.Login)

	r.Use("/ws", signup.Auth(authComponent))
	r.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})
	r.Get("/ws/user/notifications", notificationsRouter.NotificationsWebSocket())

	r.Post("/user/:id/notifications", notificationsRouter.SendNotificationToUser)

	r.Use("/docs/*", basicauth.New(basicauth.Config{
		Users: map[string]string{
			s.Config.Docs.Username: s.Config.Docs.Password,
		},
	}))
	r.Static("/docs/", s.Config.Docs.Dir)

	return nil
}
