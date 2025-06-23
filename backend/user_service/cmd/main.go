package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"

	"user_service/config"
	"user_service/routes"
)

func main() {
	// Load env
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("❌ Error loading .env file: %v", err)
	}
	// Initialize DB connection
	db, err := config.ConnectDatabase()
	if err != nil {
		log.Fatalf("❌ Database connection failed: %v", err)
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
