package tests

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"

	"video_content_management_system/backend/course_service/config"
	"video_content_management_system/backend/course_service/routes"
)

var app *fiber.App
var db *sql.DB

func TestMain(m *testing.M) {
	_ = godotenv.Load("../../.env")

	var err error
	db, err = config.ConnectDatabase()
	if err != nil {
		panic("Failed to connect to DB: " + err.Error())
	}

	app = fiber.New()
	routes.SetupCourseRoutes(app, db)

	code := m.Run()

	db.Close()
	os.Exit(code)
}

func TestGetAllCourses(t *testing.T) {
	req := httptest.NewRequest("GET", "/courses", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected 200 OK, got %d", resp.StatusCode)
	}
}

func TestCreateCourse(t *testing.T) {
	payload := map[string]interface{}{
		"title":         "Go Language Course",
		"description":   "Learn Go from scratch",
		"category_id":   1,
		"instructor_id": 1,
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest("POST", "/courses", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected 201 Created, got %d", resp.StatusCode)
	}
}
