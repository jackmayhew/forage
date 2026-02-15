package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func downloadTrack(artist, track, outputDir string) error {
	query := fmt.Sprintf("%s %s audio", artist, track)
	
	fmt.Printf("Downloading: %s - %s\n", artist, track)
	
	// Create output directory if it doesn't exist
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %v", err)
	}
	
	outputPath := filepath.Join(outputDir, "%(title)s.%(ext)s")
	
	cmd := exec.Command("yt-dlp", 
		"-x",
		"--audio-format", "mp3",
		"-o", outputPath,
		fmt.Sprintf("ytsearch1:%s", query))
	
	output, err := cmd.CombinedOutput()
	
	if err != nil {
		return fmt.Errorf("download failed: %v\n%s", err, string(output))
	}
	
	fmt.Printf("✓ Downloaded\n\n")
	return nil
}