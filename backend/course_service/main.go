package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"

	"github.com/backlog-developer/video_content_management_system/backend/course_service/routes"

	"github.com/backlog-developer/video_content_management_system/backend/course_service/config"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	app := fiber.New()

	db, err := config.ConnectDatabase()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal("Database not reachable:", err)
	}
	defer db.Close()

	// Clear naming - specific to course service
	routes.SetupCourseRoutes(app, db)

	port := os.Getenv("PORT")
	if port == "" {
		port = "4001" // Default if .env fails
	}

	log.Println("Course Service running on port", port)
	log.Fatal(app.Listen(":" + port))
}
