package api

//go:generate go run github.com/swaggo/swag/cmd/swag fmt --dir ../../../.
//go:generate go run github.com/swaggo/swag/cmd/swag init --parseDependency --dir ../../../cmd/service,service  --output ../../../docs --outputTypes json,yaml
