package main

import (
	"fmt"
	"os/exec"
	"regexp"
	"runtime"
)

var trackIDRegex = regexp.MustCompile(`track/([a-zA-Z0-9]{22})`)

func extractTrackID(url string) string {
	matches := trackIDRegex.FindStringSubmatch(url)

	if len(matches) > 1 {
		return matches[1]
	}

	return ""
}

func openFile(path string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", path)
	case "linux":
		cmd = exec.Command("xdg-open", path)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", path)
	default:
		fmt.Printf("Please open the config file manually at: %s\n", path)
		return
	}
	_ = cmd.Run()
}