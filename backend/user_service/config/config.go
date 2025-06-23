package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq" // PostgreSQL driver
)

// ConnectDatabase connects to PostgreSQL and returns the DB handle and an error if any.
func ConnectDatabase() (*sql.DB, error) {
	// Print current directory for debugging
	cwd, _ := os.Getwd()
	log.Println("📁 Current working directory:", cwd)

	// Get environment variables
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	// Debug output for env vars (except password)
	log.Println("🔧 DB Config:")
	log.Println("  DB_HOST =", dbHost)
	log.Println("  DB_PORT =", dbPort)
	log.Println("  DB_USER =", dbUser)
	log.Println("  DB_NAME =", dbName)

	psqlInfo := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName,
	)

	// Try to open DB
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Printf("❌ Failed to open database: %v\n", err)
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Try to ping DB
	if err = db.Ping(); err != nil {
		log.Printf("❌ Failed to connect to database: %v\n", err)
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("✅ Connected to PostgreSQL")
	return db, nil
}
