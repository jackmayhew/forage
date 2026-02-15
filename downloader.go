package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func downloadTrack(artist, track, outputDir string) error {
	logInfo("Downloading: %s - %s\n", artist, track)

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %v", err)
	}

	// Create filename from artist and track
	filename := fmt.Sprintf("%s - %s.mp3", artist, track)
	outputPath := filepath.Join(outputDir, filename)

	// Check if file already exists
	if _, err := os.Stat(outputPath); err == nil {
		logInfo("⊘ Already exists, skipping\n\n")
		return fmt.Errorf("skipped")
	}

	query := fmt.Sprintf("%s %s audio", artist, track)

	cmd := exec.Command("yt-dlp",
		"-x",
		"--audio-format", "mp3",
		"-o", outputPath,
		fmt.Sprintf("ytsearch1:%s", query))

	output, err := cmd.CombinedOutput()

	if err != nil {
		return fmt.Errorf("download failed: %v\n%s", err, string(output))
	}

	logInfo("✓ Downloaded\n\n")
	return nil
}