package websocket

import (
	"sidewarslobby/app/controllers"

	"github.com/antoniodipinto/ikisocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func Create(app *fiber.App) {
	app.Use(func(c *fiber.Ctx) error {
		// IsWebSocketUpgrade returns true if the client
		// requested upgrade to the WebSocket protocol.
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	ikisocket.On(ikisocket.EventDisconnect, func(ep *ikisocket.EventPayload) {
		controllers.QueueWebsocketHandleDisconnect(ep)
	})

	app.Get("/api/v1/user/queue/:token", ikisocket.New(controllers.QueueWebsocketNew))
}
