package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type DownloadJob struct {
    Artist, Title, Album, ArtURL string
}

func main() {
	godotenv.Load()

	if len(os.Args) > 1 && os.Args[1] == "config" {
		handleConfig()
		return
	}
	
	// Flags
	countFlag := flag.Int("count", 10, "Number of similar tracks to find")
	outputFlag := flag.String("output", "./foraged-tracks", "Output directory for foraged tracks")
	quietFlag := flag.Bool("quiet", false, "Quiet mode - minimal output")
	onlyFlag := flag.Bool("only", false, "Only download the provided track")
	includeSourceFlag := flag.Bool("include-source", false, "Include the provided track in the download")
	textInputFlag := flag.String("text", "", "Search for a track by 'Artist - Track'")
	flag.Parse()
	
	if *textInputFlag != "" && strings.HasPrefix(*textInputFlag, "-") {
		logAlways("Error: search text cannot start with '-' (did you forget the value for --text?)\n")
		os.Exit(1)
	}

	setQuietMode(*quietFlag)
	
	isFlagPassed := func(name string) bool {
		found := false
		flag.Visit(func(f *flag.Flag) {
			if f.Name == name { found = true }
		})
		return found
	}

	config, err := loadConfig()
	if err != nil {
		logAlways("Error: %v. Run 'forage config' to set up.\n", err)
		os.Exit(1)
	}

	spotifyClientID := config.SpotifyClientID
	spotifyClientSecret := config.SpotifyClientSecret
	lastfmAPIKey := config.LastFmAPIKey

	// Spotify access token
	token, err := getSpotifyToken(spotifyClientID, spotifyClientSecret)
	if err != nil {
		logError("Error getting Spotify token: %v\n", err)
		os.Exit(1)
	}

	if !isFlagPassed("count") && config.DefaultCount > 0 { *countFlag = config.DefaultCount }
	if !isFlagPassed("output") && config.OutputDir != "" { *outputFlag = config.OutputDir }
	if !isFlagPassed("quiet") { *quietFlag = config.QuietMode }
	if !isFlagPassed("include-source") { *includeSourceFlag = config.IncludeSource }

	setQuietMode(*quietFlag)
	
	if *countFlag > 50 {
		logAlways("Count cannot exceed 50 (Last.fm API limit)")
		os.Exit(1)
	}

	args := flag.Args()
	if len(args) < 1 && *textInputFlag == "" {
		logAlways("Usage:\n  forage <spotify-url>\n  forage --text 'Artist - Track'\n  forage config\n")
		os.Exit(1)
	}

	// Resolve source track (Spotify)
	var track *SpotifyTrack
	input := ""
	if *textInputFlag != "" {
		input = *textInputFlag
	} else if len(args) > 0 {
		input = args[0]
	}

	if strings.Contains(input, "open.spotify.com") {
		trackID := extractTrackID(input)
		track, err = getTrackInfo(token, trackID)
	} else {
		track, err = searchTrackGeneral(token, input)
	}

	if err != nil || track == nil {
		logAlways("Error: Could not find track on Spotify.\n")
		os.Exit(1)
	}

	artistName := track.Artists[0].Name
	trackName := track.Name
	logInfo("\nFound: %s - %s\n\n", artistName, trackName)

	// Find similar tracks (Last.fm)
	var similarTracks []LastFmTrack
	if !*onlyFlag {
		logInfo("Finding %d similar tracks on Last.fm...\n", *countFlag)
		similarTracks, err = getSimilarTracks(lastfmAPIKey, artistName, trackName, *countFlag)
		if err != nil {
			logError("Error getting similar tracks: %v\n", err)
			os.Exit(1)
		}
	}

	if !*onlyFlag && len(similarTracks) > 0 {
		if *quietFlag {
			logAlways("Found %d similar tracks.\n", len(similarTracks))
		} else {
			logAlways("\nFound %d similar tracks:\n\n", len(similarTracks))
			for i, t := range similarTracks {
				logAlways("%d. %s - %s\n", i+1, t.Artist.Name, t.Name)
			}
		}
	}

	// Build download queue (metadata lookup)
	var jobs []DownloadJob

	// Source track first
	if *onlyFlag || *includeSourceFlag {
		jobs = append(jobs, DownloadJob{
			Artist: artistName,
			Title:  trackName,
			Album:  track.Album.Name,
			ArtURL: track.Album.Images[0].URL,
		})
	}

	// Add similar tracks (Fetch Spotify metadata for each)
	for _, t := range similarTracks {
		meta, err := searchTrackMetadata(token, t.Artist.Name, t.Name)
		job := DownloadJob{Artist: t.Artist.Name, Title: t.Name}
		
		if err == nil && meta != nil {
			job.Album = meta.Album.Name
			if len(meta.Album.Images) > 0 {
				job.ArtURL = meta.Album.Images[0].URL
			}
		}
		jobs = append(jobs, job)
	}

	// Download
	totalToDownload := len(jobs)
	sourceSuffix := ""
	if (*includeSourceFlag || config.IncludeSource) && !*onlyFlag {
		sourceSuffix = " - including source track"
	}
	
	logAlways("\n--- Starting downloads (%d total%s) ---\n\n", totalToDownload, sourceSuffix)

	jobsChan := make(chan DownloadJob, totalToDownload)
	resultsChan := make(chan Result, totalToDownload)

	for w := 1; w <= 3; w++ {
		go worker(jobsChan, resultsChan, *outputFlag)
	}

	// Send to workers
	for _, job := range jobs {
		jobsChan <- job
	}
	close(jobsChan)

	// Collect results
	var failures []string
	successCount, skippedCount := 0, 0

	for i := 0; i < totalToDownload; i++ {
		res := <-resultsChan
		if res.Err != nil {
			if errors.Is(res.Err, ErrSkipped) {
				skippedCount++
				logInfo("[-] Skipped: %s - %s\n", res.Job.Artist, res.Job.Title)
			} else {
				failures = append(failures, fmt.Sprintf("%s - %s", res.Job.Artist, res.Job.Title))
				logInfo("[X] Failed:  %s - %s\n", res.Job.Artist, res.Job.Title)
			}
		} else {
			successCount++
			logInfo("[+] %s - %s\n", res.Job.Artist, res.Job.Title)
		}
	}

	// Summary
	logAlways("\n--- Download Summary ---\n")

	totalProcessed := successCount + skippedCount + len(failures)

	if totalProcessed == 0 {
		logAlways("No tracks processed.\n")
	} else {
		if successCount > 0 {
			trackWord := "tracks"
			if successCount == 1 { trackWord = "track" }
			logAlways("✓ Downloaded: %d %s to %s\n", successCount, trackWord, *outputFlag)
		}
		if skippedCount > 0 {
			trackWord := "tracks"
			if skippedCount == 1 { trackWord = "track" }
			logAlways("⊘ Skipped: %d %s (already exist)\n", skippedCount, trackWord)
		}
		if len(failures) > 0 {
			logAlways("\n✗ Failed downloads:\n")
			for _, track := range failures {
				logAlways("  - %s\n", track)
			}
		} else if successCount > 0 {
			logAlways("\n✓ All downloads completed successfully!\n")
		}
	}
}

type Result struct {
	Job DownloadJob
	Err error
}

func worker(jobs <-chan DownloadJob, results chan<- Result, outputDir string) {
	for job := range jobs {
		err := downloadTrack(job.Artist, job.Title, outputDir, job.Album, job.ArtURL)
		results <- Result{Job: job, Err: err}
	}
}