package routes

import (
	"os"

	"github.com/backlog-developer/video_content_management_system/backend/API_gateway_service/middleware"
	"github.com/backlog-developer/video_content_management_system/backend/API_gateway_service/utils"

	"github.com/gofiber/fiber/v2"
)

func RegisterUploadRoutes(app *fiber.App) {
	uploadServiceURL := os.Getenv("UPLOAD_SERVICE_URL")
	api := app.Group("/api/upload")

	api.Post("/", middleware.AuthMiddleware, func(c *fiber.Ctx) error {
		return utils.ForwardRequestWithBody(c, uploadServiceURL+"/api/upload/")
	})

	api.Get("/videos", middleware.AuthMiddleware, func(c *fiber.Ctx) error {
		return utils.ForwardRequest(c, uploadServiceURL+"/api/upload/videos")
	})
}
