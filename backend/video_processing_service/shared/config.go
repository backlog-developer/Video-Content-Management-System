// File: shared/config.go
package shared

import "os"

type AppConfig struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	InputDir   string
	OutputDir  string
	FFmpegPath string
}

var Config AppConfig

func LoadConfig() {
	Config = AppConfig{
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		InputDir:   os.Getenv("VIDEO_INPUT_DIR"),
		OutputDir:  os.Getenv("VIDEO_OUTPUT_DIR"),
		FFmpegPath: os.Getenv("FFMPEG_BIN"),
	}
}
