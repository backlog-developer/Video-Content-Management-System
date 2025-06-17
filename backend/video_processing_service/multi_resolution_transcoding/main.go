package multi_resolution_transcoding

import (
	"github.com/backlog-developer/video_processing_service/shared"
)

// and initializing the logger and database connection.

// main.go
// Package main initializes the video processing service, loading environment variables,

func main() {
	shared.LoadEnv()
	shared.LoadConfig()
	shared.InitLogger()
	shared.InitDB()

	shared.Info.Println("Video processing service started.")
}
