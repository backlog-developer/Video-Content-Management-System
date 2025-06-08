package routes

import (
	"github.com/gofiber/fiber/v2"
)

func Profile(c *fiber.Ctx) error {
	username := c.Locals("user").(string)
	return c.JSON(fiber.Map{
		"user":    username,
		"message": "Welcome to your profile!",
	})
}
