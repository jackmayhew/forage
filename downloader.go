package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func downloadTrack(artist, track, outputDir string) error {
	query := fmt.Sprintf("%s %s audio", artist, track)
	
	fmt.Printf("Downloading: %s - %s\n", artist, track)
	
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %v", err)
	}
	
	outputPath := filepath.Join(outputDir, "%(title)s.%(ext)s")
	archivePath := filepath.Join(outputDir, ".forage-archive.txt")
	
	cmd := exec.Command("yt-dlp", 
		"-x",
		"--audio-format", "mp3",
		"-o", outputPath,
		"--download-archive", archivePath,
		fmt.Sprintf("ytsearch1:%s", query))
	
	output, err := cmd.CombinedOutput()
	
	// Check if skipped
	if strings.Contains(string(output), "has already been recorded in the archive") {
		fmt.Printf("⊘ Already exists, skipping\n\n")
		return fmt.Errorf("skipped")
	}
	
	if err != nil {
		return fmt.Errorf("download failed: %v\n%s", err, string(output))
	}
	
	fmt.Printf("✓ Downloaded\n\n")
	return nil
}