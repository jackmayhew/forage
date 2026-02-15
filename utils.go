package main

import (
	"regexp"
)

func extractTrackID(url string) string {
	re := regexp.MustCompile(`track/([a-zA-Z0-9]{22})`)
	matches := re.FindStringSubmatch(url)

	if len(matches) > 1 {
		return matches[1]
	}

	return ""
}