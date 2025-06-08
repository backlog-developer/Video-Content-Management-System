package routes

import (
	"database/sql"
	"time"

	// Import necessary packages
	"video_content_management_system/backend/user_service/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// Load from env for security; set JWT_SECRET in .env
// JWTSecret should be set in your .env file for security
// Ensure you have JWT_SECRET defined in your .env file, e.g.:
// LoginRequest represents expected JSON payload for login
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func generateJWT(username string) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(utils.JWTSecret)
}

func Login(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var body LoginRequest
		if err := c.BodyParser(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
		}

		// Step 4: Simulate DB user check (replace with real user query later)
		// TODO: Replace this with real DB check + bcrypt password check
		if body.Username != "admin" || body.Password != "password" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
		}

		token, err := generateJWT(body.Username)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not generate token"})
		}

		return c.JSON(fiber.Map{"token": token})
	}
}
