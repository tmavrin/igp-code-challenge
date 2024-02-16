package signup

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/tmavrin/igp-code-challenge/internal/component/notifications"
	"github.com/tmavrin/igp-code-challenge/internal/types"
)

type notificationsRouter struct {
	component notifications.Provider
}

func NewNotificationsRouter(component notifications.Provider) *notificationsRouter {
	return &notificationsRouter{component: component}
}

func (r *notificationsRouter) NotificationsWebSocket() func(f *fiber.Ctx) error {
	return websocket.New(func(c *websocket.Conn) {
		defer func() {
			r.component.UnregisterClient(c)
			c.Close()
		}()

		account := c.Locals("account").(types.Account)

		// Register the client
		r.component.RegisterClient(account.ID, c)

		for {
			// keep connection open by waiting for error or close message
			_, _, err := c.ReadMessage()
			if err != nil {
				return
			}
		}

	})
}

type notificationRequest struct {
	Message string `json:"message"`
}

func (r *notificationsRouter) SendNotificationToUser(f *fiber.Ctx) error {
	accountID, err := uuid.Parse(f.Params("id"))
	if err != nil {
		return fiber.ErrBadRequest
	}

	var req notificationRequest
	err = f.BodyParser(&req)
	if err != nil {
		return fiber.ErrBadRequest
	}

	r.component.SendNotificationToUser(accountID, req.Message)
	return err
}
