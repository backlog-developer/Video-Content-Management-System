package multi_resolution_transcoding

import (
	"bytes" // Import bytes for the buffer
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/backlog-developer/video_processing_service/shared"
	"github.com/backlog-developer/video_processing_service/shared/utils"
)

func TranscodeToResolutions(inputPath, outputDir, baseFilename string) error {
	for label, size := range utils.Resolutions {
		// Construct the output file path
		// --- Important: Use filepath.Join for output path as well ---
		// This will create a path like C:\...\multi_resolution_processed\videos\original_filename_1080p.mp4
		output := filepath.Join(outputDir, fmt.Sprintf("%s_%s.mp4", baseFilename, label))
		// Ensure shared.Config.FFmpegPath is correctly populated
		ffmpegBinary := shared.Config.FFmpegPath
		if ffmpegBinary == "" {
			// Fallback if FFMPEG_BIN is not set, assuming it's in PATH
			ffmpegBinary = "ffmpeg"
			shared.Info.Println("FFMPEG_BIN env var not set, attempting to use 'ffmpeg' from PATH.")
		}

		cmdArgs := []string{"-i", inputPath, "-s", size, "-c:a", "copy", output}
		cmd := exec.Command(ffmpegBinary, cmdArgs...)

		// --- CRITICAL PART: Capture stderr and stdout ---
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out    // Capture standard output
		cmd.Stderr = &stderr // Capture standard error

		shared.Info.Printf("Attempting to execute FFmpeg command: %s %s\n", ffmpegBinary, strings.Join(cmdArgs, " "))

		err := cmd.Run()
		if err != nil {
			// Log FFmpeg's stdout/stderr when it fails
			shared.Error.Printf("FFmpeg Stdout:\n%s\n", out.String())
			shared.Error.Printf("FFmpeg Stderr:\n%s\n", stderr.String())
			return fmt.Errorf("failed to transcode to %s (%s): %w. FFmpeg process error: %s", label, size, err, stderr.String())
		}
		shared.Info.Printf("Successfully transcoded %s to %s.\n", inputPath, label)
	}
	return nil
}
