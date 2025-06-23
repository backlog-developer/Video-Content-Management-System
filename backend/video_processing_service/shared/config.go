// File: shared/config.go
package shared

import (
	"fmt"
	"os"
)

type AppConfig struct {
	DBHost             string
	DBPort             string
	DBUser             string
	DBPassword         string
	DBName             string
	FFmpegPath         string
	UploadedVideosDir  string // New: Absolute path to the upload service's storage dir
	ProcessedVideosDir string
}

var Config AppConfig

func LoadConfig() {
	Config = AppConfig{
		DBHost:             os.Getenv("DB_HOST"),
		DBPort:             os.Getenv("DB_PORT"),
		DBUser:             os.Getenv("DB_USER"),
		DBPassword:         os.Getenv("DB_PASSWORD"),
		DBName:             os.Getenv("DB_NAME"),
		FFmpegPath:         os.Getenv("FFMPEG_BIN"),
		UploadedVideosDir:  os.Getenv("UPLOADED_VIDEOS_DIR"),  // Load from .env
		ProcessedVideosDir: os.Getenv("PROCESSED_VIDEOS_DIR"), // Load from .env
	}
}

// Dsn returns the PostgreSQL connection string.
// This is the new method you need to add.
// Dsn returns the PostgreSQL connection string.
func (c *AppConfig) Dsn() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName)
}

//In shared/config.go, AppConfig is a struct (a data type), and Config is a variable of that AppConfig struct type. The LoadConfig() function populates the Config variable.
