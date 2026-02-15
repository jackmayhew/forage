package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func downloadTrack(artist, track, outputDir, album, albumArtURL string) error {
	logInfo("Downloading: %s - %s\n", artist, track)

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %v", err)
	}

	filename := fmt.Sprintf("%s - %s.mp3", artist, track)
	outputPath := filepath.Join(outputDir, filename)

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

	// Add metadata
	if err := addMetadata(outputPath, artist, track, album, albumArtURL); err != nil {
		logInfo("⚠ Downloaded but failed to add metadata\n\n")
	} else {
		logInfo("✓ Downloaded with metadata\n\n")
	}

	return nil
}