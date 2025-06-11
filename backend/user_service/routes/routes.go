package routes

import (
	"database/sql"
	"video_content_management_system/backend/user_service/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, db *sql.DB) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to VCMS Backend ARYAN HERE")
	})

	// Public routes
	app.Post("/register", Register(db))
	app.Post("/login", Login(db))

	// Protected routes
	app.Get("/profile", middleware.JWTMiddleware, Profile())
	app.Get("/me", middleware.JWTMiddleware, Me(db))
}
