package main

import (
	"fmt"
	"os/exec"
)

func downloadTrack(artist, track string) error {
	query := fmt.Sprintf("%s %s audio", artist, track)
	
	fmt.Printf("Downloading: %s - %s\n", artist, track)
	
	cmd := exec.Command("yt-dlp", 
		"-x",
		"--audio-format", "mp3",
		"-o", "%(title)s.%(ext)s",
		fmt.Sprintf("ytsearch1:%s", query))
	
	output, err := cmd.CombinedOutput()
	
	if err != nil {
		return fmt.Errorf("download failed: %v\n%s", err, string(output))
	}
	
	fmt.Printf("✓ Downloaded\n\n")
	return nil
}