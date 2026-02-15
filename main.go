package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		os.Exit(1)
	}

	spotifyClientID := os.Getenv("SPOTIFY_CLIENT_ID")
	spotifyClientSecret := os.Getenv("SPOTIFY_CLIENT_SECRET")
	lastfmAPIKey := os.Getenv("LASTFM_API_KEY")

	if spotifyClientID == "" || spotifyClientSecret == "" || lastfmAPIKey == "" {
		fmt.Println("Missing credentials in .env file")
		os.Exit(1)
	}

	// Flags
	countFlag := flag.Int("count", 10, "Number of similar tracks to find")
	outputFlag := flag.String("output", ".", "Output directory for downloaded tracks")
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("Usage: forage <spotify-url> [--count N]")
		os.Exit(1)
	}

	spotifyURL := args[0]
	trackID := extractTrackID(spotifyURL)

	if trackID == "" {
		fmt.Println("Invalid Spotify URL")
		os.Exit(1)
	}

	fmt.Printf("Track ID: %s\n", trackID)

	// Spotify access token
	token, err := getSpotifyToken(spotifyClientID, spotifyClientSecret)
	if err != nil {
		fmt.Printf("Error getting Spotify token: %v\n", err)
		os.Exit(1)
	}

	// Track info from Spotify
	track, err := getTrackInfo(token, trackID)
	if err != nil {
		fmt.Printf("Error getting track info: %v\n", err)
		os.Exit(1)
	}
	
	if len(track.Artists) == 0 {
		fmt.Println("Track not found or has no artist information")
		os.Exit(1)
	}

	artistName := track.Artists[0].Name
	trackName := track.Name

	fmt.Printf("\nFound: %s - %s\n\n", artistName, trackName)

	// Similar tracks from Last.fm
	fmt.Printf("Finding %d similar tracks on Last.fm...\n", *countFlag)
	similarTracks, err := getSimilarTracks(lastfmAPIKey, artistName, trackName, *countFlag)
	if err != nil {
		fmt.Printf("Error getting similar tracks: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\nFound %d similar tracks:\n\n", len(similarTracks))
	for i, t := range similarTracks {
		fmt.Printf("%d. %s - %s\n", i+1, t.Artist.Name, t.Name)
	}

	// Download
	fmt.Println("\n--- Starting downloads ---")
	
	var failures []string
	successCount := 0
	
	for _, t := range similarTracks {
	err := downloadTrack(t.Artist.Name, t.Name, *outputFlag)
	if err != nil {
		failures = append(failures, fmt.Sprintf("%s - %s", t.Artist.Name, t.Name))
		fmt.Printf("✗ Failed\n\n")
	} else {
		successCount++
	}
}

	// Summary
	fmt.Println("--- Download Summary ---")
	fmt.Printf("✓ Successfully downloaded: %d/%d tracks\n", successCount, len(similarTracks))
	
	if len(failures) > 0 {
		fmt.Printf("\n✗ Failed downloads:\n")
		for _, track := range failures {
			fmt.Printf("  - %s\n", track)
		}
	} else {
		fmt.Println("\nAll downloads completed successfully!")
	}
}