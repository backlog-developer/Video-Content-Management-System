package controllers

import (
	"fmt"
	"time"

	"video_content_management_system/backend/upload_service/config"
	"video_content_management_system/backend/upload_service/middleware"

	"github.com/gofiber/fiber/v2"
)

func UploadVideo(c *fiber.Ctx) error {
	userClaims := c.Locals("user").(*middleware.UserClaims)

	// Only instructors or admins can upload
	if userClaims.Role != "instructor" && userClaims.Role != "admin" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	// Get the video file
	file, err := c.FormFile("video")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "File not found"})
	}

	title := c.FormValue("title")
	description := c.FormValue("description")
	courseID := c.FormValue("course_id")

	if title == "" || courseID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Title and Course ID are required"})
	}

	// File validation
	if file.Size > 500*1024*1024 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "File size exceeds 500MB"})
	}
	if file.Header["Content-Type"][0] != "video/mp4" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Only MP4 videos allowed"})
	}

	// Save file to disk
	filePath := fmt.Sprintf("./storage/videos/%d_%s", time.Now().Unix(), file.Filename)
	if err := c.SaveFile(file, filePath); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save file"})
	}

	// Insert metadata into DB
	db := config.DB
	query := `
		INSERT INTO videos_uploads (title, description, filename, file_path, file_size, uploaded_by, course_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`

	var videoID int
	err = db.QueryRow(
		query,
		title,
		description,
		file.Filename,
		filePath,
		file.Size,
		userClaims.ID,
		courseID,
	).Scan(&videoID)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save metadata"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":  "Upload successful",
		"video_id": videoID,
		"path":     filePath,
	})
}
