package signup

import (
	"context"
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	ServiceName string `envconfig:"SERVICE"`
	Port        string `envconfig:"PORT" default:"8080"`

	DatabaseURI string `envconfig:"DB_URI" required:"true"`

	Auth struct {
		JWTKey      []byte        `envconfig:"JWT_KEY" required:"true"`
		JWTDuration time.Duration `envconfig:"JWT_DURATION" default:"24h"`
	}

	Docs struct {
		Username string `envconfig:"DOCS_USERNAME" required:"true"`
		Password string `envconfig:"DOCS_PASSWORD" required:"true"`
		Dir      string `envconfig:"DOCS_DIR" required:"true"`
	}
}

func InitConfig(ctx context.Context) (Config, error) {
	var c Config

	err := envconfig.Process("", &c)
	if err != nil {
		return c, fmt.Errorf("failed to fetch environment variables: %w", err)
	}

	return c, nil
}
