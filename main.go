package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	
	// Flags
	countFlag := flag.Int("count", 10, "Number of similar tracks to find")
	outputFlag := flag.String("output", "./foraged-tracks", "Output directory for foraged tracks")
	quietFlag := flag.Bool("quiet", false, "Quiet mode - minimal output")
	onlyFlag := flag.Bool("only", false, "Only download the provided track")
	includeSourceFlag := flag.Bool("include-source", false, "Include the provided track in the download")
	setupFlag := flag.Bool("setup", false, "Create config file template")
	flag.Parse()
	
	setQuietMode(*quietFlag)
	
	if *setupFlag {
		configPath, err := getConfigPath()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		
		if err := createConfigTemplate(); err != nil {
			fmt.Printf("Error creating config template: %v\n", err)
			os.Exit(1)
		}
		
		fmt.Printf("✓ Created config template at: %s\n", configPath)
		fmt.Println("\nAdd your API credentials:")
		fmt.Println("- Spotify: https://developer.spotify.com/dashboard")
		fmt.Println("- Last.fm: https://www.last.fm/api/account/create")
		os.Exit(0)
	}
	
	config, err := loadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		fmt.Println("\nTo set up credentials, create ~/.config/forage/config.yaml:")
		fmt.Println("spotify_client_id: your_id")
		fmt.Println("spotify_client_secret: your_secret")
		fmt.Println("lastfm_api_key: your_key")
		os.Exit(1)
	}
	
	isFlagPassed := func(name string) bool {
		found := false
		flag.Visit(func(f *flag.Flag) {
			if f.Name == name { found = true }
		})
		return found
	}

	spotifyClientID := config.SpotifyClientID
	spotifyClientSecret := config.SpotifyClientSecret
	lastfmAPIKey := config.LastFmAPIKey
	if !isFlagPassed("count") && config.DefaultCount > 0 { *countFlag = config.DefaultCount }
	if !isFlagPassed("output") && config.OutputDir != "" { *outputFlag = config.OutputDir }
	if !isFlagPassed("quiet") { *quietFlag = config.QuietMode }
	if !isFlagPassed("include-source") { *includeSourceFlag = config.IncludeSource }

	setQuietMode(*quietFlag)
	
	if *countFlag > 50 {
		fmt.Println("Count cannot exceed 50 (Last.fm API limit)")
		os.Exit(1)
	}

	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("Usage: forage <spotify-url> [--count N] [--output DIR] [--quiet]")
		os.Exit(1)
	}

	spotifyURL := args[0]
	trackID := extractTrackID(spotifyURL)

	if trackID == "" {
		fmt.Println("Invalid Spotify URL")
		os.Exit(1)
	}

	logInfo("Track ID: %s\n", trackID)

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

	logInfo("\nFound: %s - %s\n\n", artistName, trackName)
	if !*onlyFlag {
		logInfo("Finding %d similar tracks on Last.fm...\n", *countFlag)
	}
	// Similar tracks from Last.fm
	var similarTracks []LastFmTrack
	
	if !*onlyFlag {
		var err error
		similarTracks, err = getSimilarTracks(lastfmAPIKey, artistName, trackName, *countFlag)
		if err != nil {
			fmt.Printf("Error getting similar tracks: %v\n", err)
			os.Exit(1)
		}
	}

	totalToDownload := len(similarTracks)
	if *onlyFlag || *includeSourceFlag {
		totalToDownload++
	}

	if !*onlyFlag && len(similarTracks) > 0 {
		logInfo("\nFound %d similar tracks:\n\n", len(similarTracks))
		for i, t := range similarTracks {
			logInfo("%d. %s - %s\n", i+1, t.Artist.Name, t.Name)
		}
	}

	logInfo("\n--- Starting downloads (%d total) ---\n\n", totalToDownload)

	var failures []string
	successCount := 0
	skippedCount := 0
	currentIdx := 1

	if *onlyFlag || *includeSourceFlag {
		artURL := ""
		if len(track.Album.Images) > 0 { artURL = track.Album.Images[0].URL }
		
		err := downloadTrack(artistName, trackName, *outputFlag, track.Album.Name, artURL, currentIdx, totalToDownload)
		if err != nil {
			if errors.Is(err, ErrSkipped) {
				skippedCount++
			} else {
				failures = append(failures, fmt.Sprintf("%s - %s", artistName, trackName))
			}
		} else {
			successCount++
		}
		currentIdx++
	}


	if !*onlyFlag {
		for _, t := range similarTracks {
			similarTrackInfo, err := getTrackInfoBySearch(token, t.Artist.Name, t.Name)
			
			var album, albumArtURL string
			if err == nil && similarTrackInfo != nil {
				album = similarTrackInfo.Album.Name
				if len(similarTrackInfo.Album.Images) > 0 {
					albumArtURL = similarTrackInfo.Album.Images[0].URL
				}
			}

			err = downloadTrack(t.Artist.Name, t.Name, *outputFlag, album, albumArtURL, currentIdx, totalToDownload)
			if err != nil {
				if errors.Is(err, ErrSkipped) {
					skippedCount++
				} else {
					failures = append(failures, fmt.Sprintf("%s - %s", t.Artist.Name, t.Name))
				}
			} else {
				successCount++
			}
			currentIdx++
		}
	}

	// Summary
	logAlways("\n--- Download Summary ---\n")
	
	if len(failures) > 0 {
		logAlways("✓ Downloaded: %d tracks to %s\n", successCount, *outputFlag)
		if skippedCount > 0 {
			logAlways("⊘ Skipped: %d tracks (already exist)\n", skippedCount)
		}
		logAlways("\n✗ Failed downloads:\n")
		for _, track := range failures {
			logAlways("  - %s\n", track)
		}
	} else {
		if skippedCount > 0 {
			logAlways("✓ Downloaded: %d tracks to %s\n", successCount, *outputFlag)
			logAlways("⊘ Skipped: %d tracks (already exist)\n", skippedCount)
		} else {
			logAlways("✓ Downloaded %d tracks to %s\n", successCount, *outputFlag)
		}
	}
}