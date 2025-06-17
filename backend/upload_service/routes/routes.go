package routes

import (
	"github.com/gofiber/fiber/v2"

	"video_content_management_system/backend/upload_service/controllers"
	"video_content_management_system/backend/upload_service/middleware"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	upload := api.Group("/upload", middleware.JWTMiddleware())
	upload.Post("/", controllers.UploadVideo)

	// 👇 Add this GET route
	upload.Get("/videos", controllers.GetUploadedVideos)
}
