package controllers

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
)

func Profile(c *fiber.Ctx) error {
	username := c.Locals("user").(string)
	role := c.Locals("role").(string)

	return c.JSON(fiber.Map{
		"user":    username,
		"role":    role,
		"message": "Welcome to your profile!",
	})
}

func Me(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		username := c.Locals("user").(string)

		var email, role string
		var subscriptionPlanID sql.NullInt64

		err := db.QueryRow(
			"SELECT email, role, subscription_plan_id FROM users WHERE username = $1",
			username,
		).Scan(&email, &role, &subscriptionPlanID)

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get user info"})
		}

		return c.JSON(fiber.Map{
			"username":             username,
			"email":                email,
			"role":                 role,
			"subscription_plan_id": subscriptionPlanID.Int64,
		})
	}
}
