package websocket

import (
	"log"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

type WebsocketController struct {
	connections map[*websocket.Conn]bool
	connLock    sync.RWMutex
}

func NewWebsocketController() *WebsocketController {
	return &WebsocketController{
		connections: make(map[*websocket.Conn]bool),
	}
}

func (a *WebsocketController) BroadcastMessage(message interface{}) {
	a.connLock.RLock()
	defer a.connLock.RUnlock()

	for conn := range a.connections {
		if err := conn.WriteJSON(message); err != nil {
			log.Printf("Error broadcasting to connection: %v", err)
		}
	}
}

func (a *WebsocketController) WsAccess(c *fiber.Ctx) error {
	return websocket.New(func(c *websocket.Conn) {
		// Register connection
		a.connLock.Lock()
		a.connections[c] = true
		a.connLock.Unlock()

		// Cleanup on disconnect
		defer func() {
			a.connLock.Lock()
			delete(a.connections, c)
			a.connLock.Unlock()
		}()

		for {
			_, msg, err := c.ReadMessage()
			if err != nil {
				log.Println("Read error:", err)
				break
			}

			// Broadcast to all connections
			a.BroadcastMessage(msg)
		}
	})(c)
}
