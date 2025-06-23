package routes

import (
	"database/sql"

	"github.com/backlog-developer/video_content_management_system/backend/course_service/controllers"

	"github.com/gofiber/fiber/v2"
)

func SetupCourseRoutes(app *fiber.App, db *sql.DB) {
	courseGroup := app.Group("/courses")

	courseGroup.Get("/", controllers.GetAllCourses(db))
	courseGroup.Get("/:id", controllers.GetCourseByID(db))
	courseGroup.Post("/", controllers.CreateCourse(db))
	courseGroup.Put("/:id", controllers.UpdateCourse(db))
	courseGroup.Delete("/:id", controllers.DeleteCourse(db))
}
