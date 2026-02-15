package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var ErrSkipped = errors.New("skipped")

func sanitizeFilename(name string) string {
	invalid := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|"}
	result := name
	for _, char := range invalid {
		result = strings.ReplaceAll(result, char, "_")
	}
	return result
}

func downloadTrack(artist, track, outputDir, album, albumArtURL string, current, total int) error {
	logInfo("Downloading: %s - %s (%d/%d)\n", artist, track, current, total)

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %v", err)
	}

	// Sanitize artist and track names for filename
	safeArtist := sanitizeFilename(artist)
	safeTrack := sanitizeFilename(track)
	filename := fmt.Sprintf("%s - %s.mp3", safeArtist, safeTrack)
	outputPath := filepath.Join(outputDir, filename)

	if _, err := os.Stat(outputPath); err == nil {
		logInfo("⊘ Already exists, skipping\n\n")
		return ErrSkipped
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