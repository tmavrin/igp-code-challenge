package notifications

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -o ../../fake/component/notifications.go --fake-name NotificationsProvider . Provider
type Provider interface {
	SendNotificationToUser(accountID uuid.UUID, message string)
	RegisterClient(accountID uuid.UUID, ws *websocket.Conn)
	UnregisterClient(ws *websocket.Conn)
}

type component struct {
	clients         map[*websocket.Conn]uuid.UUID
	registerChan    chan wsReq
	sendMessageChan chan messageReq
	unregisterChan  chan *websocket.Conn
}

type wsReq struct {
	Socket    *websocket.Conn
	AccountID uuid.UUID
}

type messageReq struct {
	AccountID uuid.UUID
	Message   string
}

var _ Provider = (*component)(nil)

func New() *component {
	c := &component{
		clients:         make(map[*websocket.Conn]uuid.UUID),
		registerChan:    make(chan wsReq),
		sendMessageChan: make(chan messageReq),
		unregisterChan:  make(chan *websocket.Conn),
	}
	go func() {
		c.runMsgService()
	}()
	return c
}

func (c *component) runMsgService() {
	for {
		select {
		case connection := <-c.registerChan:
			c.clients[connection.Socket] = connection.AccountID

		case messageRequest := <-c.sendMessageChan:
			for connection := range c.clients {
				if c.clients[connection] == messageRequest.AccountID {
					err := connection.WriteMessage(websocket.TextMessage, []byte(messageRequest.Message))
					if err != nil {
						c.unregisterChan <- connection
						connection.WriteMessage(websocket.CloseMessage, []byte{})
						connection.Close()
					}
				}

			}

		case connection := <-c.unregisterChan:
			delete(c.clients, connection)
		}
	}
}

func (c *component) SendNotificationToUser(accountID uuid.UUID, message string) {
	c.sendMessageChan <- messageReq{AccountID: accountID, Message: message}
}

func (c *component) RegisterClient(accountID uuid.UUID, ws *websocket.Conn) {
	c.registerChan <- wsReq{AccountID: accountID, Socket: ws}
}

func (c *component) UnregisterClient(ws *websocket.Conn) {
	c.unregisterChan <- ws
}
