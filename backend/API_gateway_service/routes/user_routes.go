package routes

import (
	"os"

	"github.com/backlog-developer/video_content_management_system/backend/API_gateway_service/middleware"
	"github.com/backlog-developer/video_content_management_system/backend/API_gateway_service/utils"

	"github.com/gofiber/fiber/v2"
)

func RegisterUserRoutes(app *fiber.App) {
	userServiceURL := os.Getenv("USER_SERVICE_URL")
	api := app.Group("/api/user")

	// Public routes
	api.Post("/register", func(c *fiber.Ctx) error {
		return utils.ForwardRequestWithBody(c, userServiceURL+"/register")
	})
	api.Post("/login", func(c *fiber.Ctx) error {
		return utils.ForwardRequestWithBody(c, userServiceURL+"/login")
	})

	// Protected routes
	api.Get("/profile", middleware.AuthMiddleware, func(c *fiber.Ctx) error {
		return utils.ForwardRequest(c, userServiceURL+"/profile")
	})
	api.Get("/me", middleware.AuthMiddleware, func(c *fiber.Ctx) error {
		return utils.ForwardRequest(c, userServiceURL+"/me")
	})
}
