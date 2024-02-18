package service

import (
	"errors"
	"net/mail"

	"github.com/gofiber/fiber/v2"
	"github.com/tmavrin/igp-code-challenge/internal/component/auth"
	"github.com/tmavrin/igp-code-challenge/internal/types"
)

type authRouter struct {
	component auth.Provider
}

func NewAuthRouter(component auth.Provider) *authRouter {
	return &authRouter{component: component}
}

// Register godoc
//
//	@Summary		Register
//	@Description	Register
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			request		body	types.AuthCredentials	true	"body"
//	@Success		200
//	@Failure		409
//	@Failure		500
//	@Router			/auth/register [POST]
func (a *authRouter) Register(f *fiber.Ctx) error {
	var credentials types.AuthCredentials
	err := f.BodyParser(&credentials)
	if err != nil {
		return err
	}
	_, err = mail.ParseAddress(credentials.Email)
	if err != nil {
		return fiber.ErrBadRequest
	}

	if len(credentials.Password) < 6 {
		return fiber.ErrBadRequest
	}

	err = a.component.Create(f.Context(), credentials)
	if errors.Is(err, auth.ErrorUserExists) {
		return fiber.ErrConflict
	}

	return err
}

type loginResponse struct {
	types.Account
	Token string `json:"token"`
}

// Login godoc
//
//	@Summary		Login
//	@Description	Login
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			request		body	types.AuthCredentials	true	"body"
//	@Success		200	{object}	loginResponse
//	@Failure		403
//	@Failure		500
//	@Router			/auth/login [POST]
func (a *authRouter) Login(f *fiber.Ctx) error {
	var credentials types.AuthCredentials
	err := f.BodyParser(&credentials)
	if err != nil {
		return err
	}

	account, token, err := a.component.Login(f.Context(), credentials)
	if errors.Is(err, auth.ErrorInvalidCredentials) {
		return fiber.ErrUnauthorized
	} else if err != nil {
		return err
	}

	return f.JSON(loginResponse{Account: account, Token: token})
}
