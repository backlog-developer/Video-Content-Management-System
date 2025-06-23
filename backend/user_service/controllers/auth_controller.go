package controllers

import (
	"database/sql"
	"time"

	"user_service/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

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

func generateJWT(username, role string) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"role":     role,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(utils.GenerateJWT)
}

// Register handles user registration
// It expects a JSON payload with "username", "email", "password", and optional "role"
// It checks if the username or email already exists, hashes the password, and inserts the user into the database
// It returns a 201 status on success or appropriate error messages
func Register(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var body RegisterRequest
		if err := c.BodyParser(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
		}
		// Step 1: Validate input
		// Ensure username, email, and password are provided
		if body.Username == "" || body.Email == "" || body.Password == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Username, Email, and Password are required",
			})
		}

		// Step 2: Check if username or email already exists
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

		// Step 3: Hash the password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to hash password"})
		}
		// Step 4: Set default role if empty
		// If role is not provided, default to "user"
		role := body.Role
		if role == "" {
			role = "user"
		}

		// Step 5: Insert into database
		_, err = db.Exec(`
            INSERT INTO users (username, email, password_hash, role, created_at)
            VALUES ($1, $2, $3, $4, $5)
        `, body.Username, body.Email, string(hashedPassword), role, time.Now())

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create user: " + err.Error()})
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "User registered successfully"})
	}
}

// Login handles user login and JWT generation
// It expects a JSON payload with "username" and "password"
func Login(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var body LoginRequest
		if err := c.BodyParser(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
		}
		// Step 1: Query user by username
		var id int
		var username, passwordHash, role string

		err := db.QueryRow(
			"SELECT id, username, password_hash, role FROM users WHERE username = $1 OR email = $1",
			body.Username,
		).Scan(&id, &username, &passwordHash, &role)

		if err != nil {
			if err == sql.ErrNoRows {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
		}

		// Step 2: Compare password hash
		err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(body.Password))
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
		}

		// Step 3: Generate JWT
		token, err := generateJWT(username, role)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not generate token"})
		}
		// Step 4: Return token
		return c.JSON(fiber.Map{"token": token})
	}
}

// Note: Ensure you have the JWT_SECRET environment variable set in your .env file
// utils/jwt.go
//Summary of flow:
//✅ User sends /login → POST → { "username": "...", "password": "..." }
//✅ Your server queries the users table → gets password_hash
//✅ Compares the provided password with stored hash
//✅ If correct → JWT token is generated → returned
//✅ If incorrect → error 401 Unauthorized //
