// File: upload_service/main.go
package main

import (
	"log"
	"os"

	"github.com/backlog-developer/video_content_management_system/backend/upload_service/config"
	"github.com/backlog-developer/video_content_management_system/backend/upload_service/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	config.InitDB() //

	app := fiber.New(fiber.Config{
		BodyLimit: 500 * 1024 * 1024, // 500 MB  limit
	})
	routes.RegisterUploadRoutes(app)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Println("🚀 Server running on port", port)
	log.Fatal(app.Listen(":" + port))
}
