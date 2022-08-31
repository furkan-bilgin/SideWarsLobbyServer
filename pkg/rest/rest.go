package rest

import (
	"github.com/gofiber/fiber/v2"
)

// Create new REST API serveer
func Create() *fiber.App {
	app := fiber.New()

	/*	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendStatus(404)
	})*/

	return app
}
