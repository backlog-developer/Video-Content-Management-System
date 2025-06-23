package test

import (
	"testing"

	"github.com/backlog-developer/video_processing_service/multi_resolution_transcoding"
	"github.com/backlog-developer/video_processing_service/shared"
	"github.com/backlog-developer/video_processing_service/shared/utils"
)

func TestTranscoding(t *testing.T) {
	shared.LoadEnv()
	shared.LoadConfig()
	shared.InitLogger()

	inputPath := shared.Config.UploadedVideosDir + "/sample_video.mp4" // ⚠️ Ensure this file exists
	baseFilename := "sample_video"

	// Ensure output directory exists
	err := utils.EnsureDir(shared.Config.ProcessedVideosDir)
	if err != nil {
		t.Fatalf("Failed to create output directory: %v", err)
	}

	// Perform transcoding
	err = multi_resolution_transcoding.TranscodeToResolutions(inputPath, shared.Config.ProcessedVideosDir, baseFilename)
	if err != nil {
		t.Fatalf("Transcoding failed: %v", err)
	}

	t.Log("✅ Transcoding completed successfully.")
}
