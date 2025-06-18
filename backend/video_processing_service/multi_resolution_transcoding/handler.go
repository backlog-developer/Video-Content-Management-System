package multi_resolution_transcoding

import (
	"github.com/backlog-developer/video_processing_service/shared"
	"github.com/backlog-developer/video_processing_service/shared/utils"
)

func StartTranscodingJob(inputPath string, baseFilename string) {
	outputDir := shared.Config.OutputDir
	err := utils.EnsureDir(outputDir)
	if err != nil {
		shared.Error.Println("Error creating output dir:", err)
		return
	}

	err = TranscodeToResolutions(inputPath, outputDir, baseFilename)
	if err != nil {
		shared.Error.Println("Transcoding failed:", err)
		return
	}

	shared.Info.Println("✅ Transcoding completed for:", inputPath)
}
