package controllers

import (
	"video_content_management_system/backend/upload_service/config"
	"video_content_management_system/backend/upload_service/middleware"

	"github.com/gofiber/fiber/v2"
)

func GetUploadedVideos(c *fiber.Ctx) error {
	user := c.Locals("user").(*middleware.UserClaims)

	query := "SELECT id, title, filename, file_path, uploaded_by FROM videos_uploads"
	args := []interface{}{}

	if user.Role == "instructor" {
		query += " WHERE uploaded_by = $1"
		args = append(args, user.ID)
	} else if user.Role == "student" {
		// Optional: Filter by course access later
	}

	rows, err := config.DB.Query(query, args...)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "DB error"})
	}
	defer rows.Close()

	var videos []map[string]interface{}
	for rows.Next() {
		var id int
		var title, filename, path string
		var uploadedBy int
		rows.Scan(&id, &title, &filename, &path, &uploadedBy)
		videos = append(videos, fiber.Map{
			"id":       id,
			"title":    title,
			"filename": filename,
			"path":     path,
		})
	}

	return c.JSON(videos)
}
