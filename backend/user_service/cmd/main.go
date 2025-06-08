package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"

	"video_content_management_system/backend/user_service/config"
	"video_content_management_system/backend/user_service/routes"
)

func main() {
	// Load env
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize DB connection
	db, err := config.ConnectDatabase()
	if err != nil {
		log.Fatalf("Database connection error: %v", err)
	}
	defer db.Close()

	// Init Fiber app
	app := fiber.New()

	// Setup routes
	routes.SetupRoutes(app, db)

	// Get port or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server running on port %s...", port)
	log.Fatal(app.Listen(":" + port))
}
