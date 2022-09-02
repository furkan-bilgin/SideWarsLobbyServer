package utils

import "github.com/gofiber/fiber/v2"

func RESTError(c *fiber.Ctx, message string) error {
	return c.JSON(fiber.Map{
		"error":   true,
		"message": message,
	})
}
