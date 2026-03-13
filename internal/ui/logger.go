package ui

import "fmt"

var quietMode bool

func SetQuietMode(quiet bool) {
	quietMode = quiet
}

func LogInfo(format string, args ...interface{}) {
	if !quietMode {
		fmt.Printf(format, args...)
	}
}

func LogError(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

func LogAlways(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}