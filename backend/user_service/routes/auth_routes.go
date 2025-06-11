package routes

import (
	"database/sql"
	"time"

	"video_content_management_system/backend/user_service/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Structs for request parsing
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

// Public Route: /register
func Register(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var body RegisterRequest
		if err := c.BodyParser(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
		}

		if body.Username == "" || body.Email == "" || body.Password == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Username, Email, and Password are required",
			})
		}

		var exists bool
		err := db.QueryRow(`
			SELECT EXISTS (SELECT 1 FROM users WHERE username=$1 OR email=$2)
		`, body.Username, body.Email).Scan(&exists)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
		}
		if exists {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Username or Email already exists"})
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to hash password"})
		}

		role := body.Role
		if role == "" {
			role = "user"
		}

		_, err = db.Exec(`
			INSERT INTO users (username, email, password_hash, role, created_at)
			VALUES ($1, $2, $3, $4, $5)
		`, body.Username, body.Email, string(hashedPassword), role, time.Now())

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create user"})
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "User registered successfully"})
	}
}

// Public Route: /login
func Login(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var body LoginRequest
		if err := c.BodyParser(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
		}

		var id int
		var username, passwordHash, role string
		err := db.QueryRow(
			"SELECT id, username, password_hash, role FROM users WHERE username = $1 OR email = $1",
			body.Username,
		).Scan(&id, &username, &passwordHash, &role)

		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
		}

		if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(body.Password)); err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
		}

		// Generate JWT
		claims := jwt.MapClaims{
			"username": username,
			"role":     role,
			"exp":      time.Now().Add(72 * time.Hour).Unix(),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(utils.JWTSecret)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not generate token"})
		}

		return c.JSON(fiber.Map{"token": tokenString})
	}
}

// Protected Route: /me
func Me(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		username := c.Locals("user").(string)
		var email, role string
		var subscriptionPlanID sql.NullInt64

		err := db.QueryRow(
			"SELECT email, role, subscription_plan_id FROM users WHERE username = $1",
			username,
		).Scan(&email, &role, &subscriptionPlanID)

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get user info"})
		}

		return c.JSON(fiber.Map{
			"username":             username,
			"email":                email,
			"role":                 role,
			"subscription_plan_id": subscriptionPlanID.Int64,
		})
	}
}

// Protected Route: /profile
func Profile() fiber.Handler {
	return func(c *fiber.Ctx) error {
		username := c.Locals("user").(string)
		role := c.Locals("role").(string)

		return c.JSON(fiber.Map{
			"user":    username,
			"role":    role,
			"message": "Welcome to your profile!",
		})
	}
}
