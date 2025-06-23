// File: multi_resolution_transcoding/main.go
package main

import (
	"database/sql"
	"path/filepath"
	"time"

	"github.com/backlog-developer/video_processing_service/multi_resolution_transcoding"
	"github.com/backlog-developer/video_processing_service/shared"
	"github.com/backlog-developer/video_processing_service/shared/utils"

	_ "github.com/lib/pq"
)

func main() {
	shared.LoadEnv()
	shared.LoadConfig()
	shared.InitLogger()

	// --- Debug prints for new config values ---
	shared.Info.Printf("FFMPEG_BIN configured as: %s\n", shared.Config.FFmpegPath)
	shared.Info.Printf("UPLOADED_VIDEOS_DIR configured as: %s\n", shared.Config.UploadedVideosDir)
	shared.Info.Printf("PROCESSED_VIDEOS_DIR configured as: %s\n", shared.Config.ProcessedVideosDir)
	// --- End debug prints ---

	// Corrected line: Call the Dsn() method on shared.Config
	db, err := sql.Open("postgres", shared.Config.Dsn())
	if err != nil {
		shared.Error.Fatalf("❌ Failed to connect to DB: %v", err)
	}
	defer db.Close()
	// Check DB connection

	shared.Info.Println("✅ Connected to database successfully.")
	shared.Info.Println("🔁 Starting auto-transcoding loop...")

	for {
		// 1. Fetch unprocessed videos
		// Ensure your DB query fetches `filename` (which seems to be the actual filename)
		rows, err := db.Query("SELECT id, file_path, filename FROM videos_uploads WHERE upload_status = 'pending'")
		if err != nil {
			shared.Error.Println("❌ DB query error:", err)
			time.Sleep(10 * time.Second)
			continue
		}

		for rows.Next() {
			var id int
			var path, filename string
			if err := rows.Scan(&id, &path, &filename); err != nil {
				shared.Error.Println("❌ Row scan error:", err)
				continue
			}

			// --- CRITICAL CHANGE: Construct absolute inputPath ---
			// Combine the absolute base directory with the filename
			inputPath := filepath.Join(shared.Config.UploadedVideosDir, filename)
			// Ensure the path is clean (removes redundant separators, etc.)
			inputPath = filepath.Clean(inputPath)
			// --- END CRITICAL CHANGE ---

			shared.Info.Printf("🎬 Transcoding video ID %d: %s (Absolute Input Path: %s)\n", id, filename, inputPath)

			// Ensure the output directory exists
			// This `utils.EnsureDir` should create the full path including subdirectories
			if err := utils.EnsureDir(shared.Config.ProcessedVideosDir); err != nil {
				shared.Error.Println("❌ Failed to ensure output dir:", err)
				continue
			}

			// Ensure output dir exists
			if err := utils.EnsureDir(shared.Config.ProcessedVideosDir); err != nil {
				shared.Error.Println("❌ Failed to ensure output dir:", err)
				continue
			}

			// Run transcoding
			err = multi_resolution_transcoding.TranscodeToResolutions(inputPath, shared.Config.ProcessedVideosDir, filename)
			if err != nil {
				shared.Error.Printf("❌ Transcoding failed for video ID %d: %v\n", id, err)
				continue
			}

			// Mark video as processed
			_, err = db.Exec("UPDATE videos_uploads SET upload_status = 'processed', updated_at = CURRENT_TIMESTAMP WHERE id = $1", id)
			if err != nil {
				shared.Error.Printf("❌ Failed to update video ID %d status: %v\n", id, err)
			} else {
				shared.Info.Printf("✅ Video ID %d marked as processed.\n", id)
			}
		}
		rows.Close()

		// Wait before next iteration
		time.Sleep(10 * time.Second)
	}
}
