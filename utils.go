package main

import (
	"regexp"
)

var trackIDRegex = regexp.MustCompile(`track/([a-zA-Z0-9]{22})`)

func extractTrackID(url string) string {
	matches := trackIDRegex.FindStringSubmatch(url)

	if len(matches) > 1 {
		return matches[1]
	}

	return ""
}