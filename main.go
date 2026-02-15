package main

import (
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

	if len(os.Args) < 2 {
		fmt.Println("Usage: music-finder <spotify-url>")
		os.Exit(1)
	}

	spotifyURL := os.Args[1]
	trackID := extractTrackID(spotifyURL)

	if trackID == "" {
		fmt.Println("Invalid Spotify URL")
		os.Exit(1)
	}

	fmt.Printf("Track ID: %s\n", trackID)

	// Get Spotify access token
	token, err := getSpotifyToken(spotifyClientID, spotifyClientSecret)
	if err != nil {
		fmt.Printf("Error getting Spotify token: %v\n", err)
		os.Exit(1)
	}

	// Get track info from Spotify
	track, err := getTrackInfo(token, trackID)
	if err != nil {
		fmt.Printf("Error getting track info: %v\n", err)
		os.Exit(1)
	}

	artistName := track.Artists[0].Name
	trackName := track.Name

	fmt.Printf("\nFound: %s - %s\n\n", artistName, trackName)

	// Get similar tracks from Last.fm
	fmt.Println("Finding similar tracks on Last.fm...")
	similarTracks, err := getSimilarTracks(lastfmAPIKey, artistName, trackName, 10)
	if err != nil {
		fmt.Printf("Error getting similar tracks: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\nFound %d similar tracks:\n\n", len(similarTracks))
	for i, t := range similarTracks {
		fmt.Printf("%d. %s - %s\n", i+1, t.Artist.Name, t.Name)
	}
}