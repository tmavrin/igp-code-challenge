package api

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/tmavrin/igp-code-challenge/internal/http/api/service"
)

type signupServer struct {
	Resource service.Resource
	Config   service.Config
	fiberApp *fiber.App
}

func NewServer(ctx context.Context) (*signupServer, error) {
	var (
		s   signupServer
		err error
	)

	s.Config, err = service.InitConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize resources: %w", err)
	}

	s.Resource, err = service.InitResource(ctx, s.Config)
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

	if s.Config.TLS.Cert != "" && s.Config.TLS.Key != "" {
		cert, err := tls.LoadX509KeyPair(s.Config.TLS.Cert, s.Config.TLS.Key)
		if err != nil {
			return err
		}
		return s.fiberApp.ListenTLSWithCertificate(":"+s.Config.Port, cert)
	} else {
		return s.fiberApp.Listen(":" + s.Config.Port)
	}
}

func (s *signupServer) Close(ctx context.Context, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return errors.Join(
		s.fiberApp.ShutdownWithContext(ctx),
		s.Resource.Close(),
	)
}
