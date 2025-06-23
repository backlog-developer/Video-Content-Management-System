package routes

import (
	"github.com/backlog-developer/video_content_management_system/backend/upload_service/controllers"
	"github.com/backlog-developer/video_content_management_system/backend/upload_service/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterUploadRoutes(app *fiber.App) {
	api := app.Group("/api/upload", middleware.JWTMiddleware())
	api.Post("/", controllers.UploadVideo)

}
