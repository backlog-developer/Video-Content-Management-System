package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"

	"github.com/backlog-developer/video_content_management_system/backend/API_gateway_service/routes"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(".env"); err != nil {
		log.Println("⚠️  .env file not found, using system environment variables.")
	}

	// Create a new Fiber app
	app := fiber.New()

	// Register all routes
	routes.RegisterUserRoutes(app)
	routes.RegisterCourseRoutes(app)
	routes.RegisterUploadRoutes(app)

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 API Gateway is running on http://localhost:%s\n", port)
	log.Fatal(app.Listen(":" + port))
}
