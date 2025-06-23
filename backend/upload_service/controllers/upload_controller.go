package controllers

import (
	"fmt"
	"strconv"

	"github.com/backlog-developer/video_content_management_system/backend/upload_service/config"
	"github.com/backlog-developer/video_content_management_system/backend/upload_service/middleware"
	"github.com/backlog-developer/video_content_management_system/backend/upload_service/storage"

	"github.com/gofiber/fiber/v2"
)

func UploadVideo(c *fiber.Ctx) error {
	userClaims := c.Locals("user").(*middleware.UserClaims)
	fmt.Printf("🔍 Uploading as user ID: %d (Role: %s)\n", userClaims.UserID, userClaims.Role)

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
	courseIDStr := c.FormValue("course_id")

	// Validate required fields
	if title == "" || courseIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Title and Course ID are required"})
	}

	courseID, err := strconv.Atoi(courseIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid course ID"})
	}

	// File validation
	if file.Size > 500*1024*1024 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "File size exceeds 500MB"})
	}
	if file.Header["Content-Type"][0] != "video/mp4" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Only MP4 videos allowed"})
	}

	// Save file using local storage module
	filePath, err := storage.SaveToLocal(c, file)
	if err != nil {
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
		userClaims.UserID,
		courseID,
	).Scan(&videoID)

	if err != nil {
		fmt.Println("❌ SQL Error inserting metadata:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save metadata"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":  "Upload successful",
		"video_id": videoID,
		"path":     filePath,
	})
}
