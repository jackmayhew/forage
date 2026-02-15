package main

import "fmt"

var quietMode bool

func setQuietMode(quiet bool) {
	quietMode = quiet
}

func logInfo(format string, args ...interface{}) {
	if !quietMode {
		fmt.Printf(format, args...)
	}
}

func logError(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

func logAlways(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}