package api

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/tmavrin/igp-code-challenge/internal/http/api/signup"
)

type signupServer struct {
	Resource signup.Resource
	Config   signup.Config
	fiberApp *fiber.App
}

func NewSignupServer(ctx context.Context) (*signupServer, error) {
	var (
		s   signupServer
		err error
	)

	s.Config, err = signup.InitConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize resources: %w", err)
	}

	s.Resource, err = signup.InitResource(ctx, s.Config)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize resources: %w", err)
	}

	s.fiberApp = fiber.New(fiber.Config{AppName: s.Config.ServiceName})
	s.fiberApp.Use(recover.New())

	return &s, nil
}

func (s *signupServer) ListenAndServe() error {
	s.fiberApp.Use(cors.New())

	s.routes()
	return s.fiberApp.Listen(":" + s.Config.Port)
}

func (s *signupServer) Close(ctx context.Context, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return errors.Join(
		s.fiberApp.ShutdownWithContext(ctx),
		s.Resource.Close(),
	)
}
