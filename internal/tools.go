//go:build tools
// +build tools

package tools

// needed for go generate ./... to generate mocks and docs in one command
import (
	_ "github.com/maxbrunsfeld/counterfeiter/v6"
	_ "github.com/swaggo/swag/cmd/swag"
)
