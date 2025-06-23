package routes

import (
	"database/sql"
	"user_service/models"
	"user_service/utils"

	"github.com/gofiber/fiber/v2"
)

// RegisterUser handles user registration using raw SQL (*sql.DB)
func RegisterUser(c *fiber.Ctx, db *sql.DB) error {
	var user models.User

	// Parse incoming JSON into user struct
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	// Check for duplicate email
	var exists int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE email = $1", user.Email).Scan(&exists)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error checking for duplicate email",
		})
	}
	if exists > 0 {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "Email already exists",
		})
	}

	// Insert user into DB
	query := `INSERT INTO users (username, email, password, role) VALUES ($1, $2, $3, $4) RETURNING id`
	// Use the correct password field; update 'HashedPassword' if your struct uses a different name
	err = db.QueryRow(query, user.Username, user.Email, user.PasswordHash, user.Role).Scan(&user.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to register user",
		})
	}

	// Generate JWT
	token, err := utils.GenerateJWT(int(user.ID), user.Role)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}

	// Return success
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User registered successfully",
		"token":   token,
	})
}
