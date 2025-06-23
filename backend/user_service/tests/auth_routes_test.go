package tests

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"

	"user_service/config"
	"user_service/routes"
)

var app *fiber.App
var db *sql.DB
var jwtToken string

// Constants → easy to update if needed
var testUsername = "testuser123"
var testEmail = "testuser123@example.com"
var testPassword = "testpassword"

// Setup runs before all tests
func TestMain(m *testing.M) {
	err := godotenv.Load("../.env") // Go up 1 level from tests/
	if err != nil {
		err = godotenv.Load(".env") // fallback to local
		if err != nil {
			fmt.Println("Warning: .env file not loaded", err)
		}
	}

	// Initialize database connection
	db, err = config.ConnectDatabase()
	if err != nil {
		panic(err)
	}

	app = fiber.New()
	routes.SetupRoutes(app, db)

	// Run tests
	code := m.Run()

	// Clean up test user after all tests
	CleanupTestUser()

	db.Close()
	os.Exit(code)
}

// Helper → ensure test user exists
func SetupTestUser(t *testing.T) {
	payload := map[string]string{
		"username": testUsername,
		"email":    testEmail,
		"password": testPassword,
		"role":     "user",
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest("POST", "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)

	if resp.StatusCode == http.StatusCreated {
		t.Log("Test user registered successfully")
	} else if resp.StatusCode == http.StatusBadRequest {
		t.Log("Test user already exists (OK)")
	} else {
		t.Fatalf("Unexpected response: %d", resp.StatusCode)
	}
}

// Helper → delete test user after tests
func CleanupTestUser() {
	_, err := db.Exec(`DELETE FROM users WHERE username = $1 OR email = $2`, testUsername, testEmail)
	if err != nil {
		// You can log this but don't panic here
		println("Warning: Failed to clean up test user:", err.Error())
	} else {
		println("Test user cleaned up")
	}
}

// Test Register Duplicate User
func TestRegisterDuplicate(t *testing.T) {
	SetupTestUser(t)

	payload := map[string]string{
		"username": testUsername,
		"email":    testEmail,
		"password": testPassword,
		"role":     "user",
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest("POST", "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", resp.StatusCode)
	}
}

// Test Login Success → Save JWT
func TestLoginSuccess(t *testing.T) {
	payload := map[string]string{
		"username": testUsername,
		"password": testPassword,
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", resp.StatusCode)
	}

	var result map[string]string
	json.NewDecoder(resp.Body).Decode(&result)

	token, exists := result["token"]
	if !exists || token == "" {
		t.Fatalf("Expected JWT token, got none")
	}

	// Save token globally
	jwtToken = token
}

// Test Login Invalid Password
func TestLoginInvalidPassword(t *testing.T) {
	payload := map[string]string{
		"username": testUsername,
		"password": "wrongpassword",
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", resp.StatusCode)
	}
}

// Test Profile with Valid JWT
func TestProfileWithValidJWT(t *testing.T) {
	if jwtToken == "" {
		t.Fatal("JWT token not set — run TestLoginSuccess first")
	}

	req := httptest.NewRequest("GET", "/profile", nil)
	req.Header.Set("Authorization", "Bearer "+jwtToken)

	resp, _ := app.Test(req)

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	if result["user"] != testUsername {
		t.Errorf("Expected username '%s', got %v", testUsername, result["user"])
	}
}

// Test Profile with Invalid JWT
func TestProfileInvalidJWT(t *testing.T) {
	req := httptest.NewRequest("GET", "/profile", nil)
	req.Header.Set("Authorization", "Bearer invalid.token.here")

	resp, _ := app.Test(req)

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", resp.StatusCode)
	}
}

// Test Me with Valid JWT
func TestMeWithValidJWT(t *testing.T) {
	if jwtToken == "" {
		t.Fatal("JWT token not set — run TestLoginSuccess first")
	}

	req := httptest.NewRequest("GET", "/me", nil)
	req.Header.Set("Authorization", "Bearer "+jwtToken)

	resp, _ := app.Test(req)

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	if result["username"] != testUsername {
		t.Errorf("Expected username '%s', got %v", testUsername, result["username"])
	}
}

// Test Me with Invalid JWT
func TestMeWithInvalidJWT(t *testing.T) {
	req := httptest.NewRequest("GET", "/me", nil)
	req.Header.Set("Authorization", "Bearer invalid.token.here")

	resp, _ := app.Test(req)

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", resp.StatusCode)
	}
}

//Recap of flow:
//✅ TestRegisterSuccess → first register the user → should succeed
//✅ TestRegisterDuplicate → register same user → should give 400
//✅ TestLoginSuccess → login → get valid JWT token → save in jwtToken
//✅ TestLoginInvalidPassword → wrong password → should give 401
//✅ TestProfileWithValidJWT → call /profile with valid JWT → should get profile info
//✅ TestProfileInvalidJWT → call /profile with garbage JWT → should get 401
// ✅ TestMeWithValidJWT → call /me with valid JWT → should get user info
// ✅ TestMeWithInvalidJWT → call /me with garbage JWT → should get 401
