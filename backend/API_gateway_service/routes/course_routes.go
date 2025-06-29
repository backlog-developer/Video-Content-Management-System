package routes

import (
	"os"

	"github.com/backlog-developer/video_content_management_system/backend/API_gateway_service/middleware"
	"github.com/backlog-developer/video_content_management_system/backend/API_gateway_service/utils"

	"github.com/gofiber/fiber/v2"
)

func RegisterCourseRoutes(app *fiber.App) {
	courseServiceURL := os.Getenv("COURSE_SERVICE_URL")
	api := app.Group("/api/course")

	// Public
	api.Get("/", func(c *fiber.Ctx) error {
		return utils.ForwardRequest(c, courseServiceURL+"/courses")
	})
	api.Get("/:id", func(c *fiber.Ctx) error {
		return utils.ForwardRequest(c, courseServiceURL+"/courses/"+c.Params("id"))
	})

	// Protected
	api.Post("/", middleware.AuthMiddleware, func(c *fiber.Ctx) error {
		return utils.ForwardRequestWithBody(c, courseServiceURL+"/courses")
	})
	api.Put("/:id", middleware.AuthMiddleware, func(c *fiber.Ctx) error {
		return utils.ForwardRequestWithBody(c, courseServiceURL+"/courses/"+c.Params("id"))
	})
	api.Delete("/:id", middleware.AuthMiddleware, func(c *fiber.Ctx) error {
		return utils.ForwardRequest(c, courseServiceURL+"/courses/"+c.Params("id"))
	})
}
