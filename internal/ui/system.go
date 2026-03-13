package ui

import (
	"fmt"
	"os/exec"
	"runtime"
)

func OpenFile(path string) {
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