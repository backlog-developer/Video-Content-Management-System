package multi_resolution_transcoding

import (
	"fmt"
	"os/exec"

	"github.com/backlog-developer/video_processing_service/shared/utils"
)

func TranscodeToResolutions(inputPath, outputDir, baseFilename string) error {
	for label, size := range utils.Resolutions {
		output := fmt.Sprintf("%s/%s_%s.mp4", outputDir, baseFilename, label)
		cmd := exec.Command("ffmpeg", "-i", inputPath, "-s", size, "-c:a", "copy", output)
		err := cmd.Run()
		if err != nil {
			return fmt.Errorf("failed to transcode to %s: %w", label, err)
		}
	}
	return nil
}
