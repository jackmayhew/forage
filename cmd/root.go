package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"forage/internal/config"
	"forage/internal/downloader"
	"forage/internal/lastfm"
	"forage/internal/spotify"
	"forage/internal/ui"

	"github.com/spf13/cobra"
)

var (
	count         int
	outputDir     string
	quiet         bool
	only          bool
	includeSource bool
	textInput     string
)

var rootCmd = &cobra.Command{
	Use:   "forage [spotify-url]",
	Short: "Find and download similar music from Spotify",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		if err != nil {
			ui.LogAlways("Error: %v. Run 'forage config' to set up.\n", err)
			os.Exit(1)
		}

		// Flag overrides config
		if !cmd.Flags().Changed("count") { count = cfg.DefaultCount }
		if !cmd.Flags().Changed("output") { outputDir = cfg.OutputDir }
		if !cmd.Flags().Changed("include-source") { includeSource = cfg.IncludeSource }
		if !cmd.Flags().Changed("quiet") { quiet = cfg.QuietMode }

		ui.SetQuietMode(quiet)

		if count > 50 {
			ui.LogAlways("Count cannot exceed 50\n")
			os.Exit(1)
		}

		token, err := spotify.GetToken(cfg.SpotifyClientID, cfg.SpotifyClientSecret)
		if err != nil {
			ui.LogError("Error getting Spotify token: %v\n", err)
			os.Exit(1)
		}

		var track *spotify.Track
		input := textInput
		if len(args) > 0 {
			input = args[0]
		}

		if input == "" {
			cmd.Help()
			return
		}

		if strings.Contains(input, "open.spotify.com") {
			trackID := spotify.ExtractTrackID(input)
			track, err = spotify.GetTrackInfo(token, trackID)
		} else {
			track, err = spotify.Search(token, input)
		}

		if err != nil || track == nil {
			ui.LogAlways("Error: Could not find track on Spotify.\n")
			os.Exit(1)
		}

		artistName := track.Artists[0].Name
		trackName := track.Name
		ui.LogInfo("\nFound: %s - %s\n\n", artistName, trackName)

		var similarTracks []lastfm.Track
		if !only {
			ui.LogInfo("Finding %d similar tracks on Last.fm...\n", count)
			similarTracks, err = lastfm.GetSimilarTracks(cfg.LastFmAPIKey, artistName, trackName, count)
			if err != nil {
				ui.LogError("Error getting similar tracks: %v\n", err)
				os.Exit(1)
			}
		}

		if !only && len(similarTracks) > 0 {
			ui.LogAlways("\nFound %d similar tracks:\n\n", len(similarTracks))
			if !quiet {
				for i, t := range similarTracks {
					ui.LogAlways("%d. %s - %s\n", i+1, t.Artist.Name, t.Name)
				}
			}
		}

		runDownloads(token, track, similarTracks)
	},
}

type DownloadJob struct {
	Artist, Title, Album, ArtURL string
}

type Result struct {
	Job DownloadJob
	Err error
}

func runDownloads(token string, source *spotify.Track, similar []lastfm.Track) {
	var jobs []DownloadJob

	if only || includeSource {
		jobs = append(jobs, DownloadJob{
			Artist: source.Artists[0].Name,
			Title:  source.Name,
			Album:  source.Album.Name,
			ArtURL: source.Album.Images[0].URL,
		})
	}

	for _, t := range similar {
		meta, err := spotify.SearchMetadata(token, t.Artist.Name, t.Name)
		job := DownloadJob{Artist: t.Artist.Name, Title: t.Name}
		if err == nil && meta != nil {
			job.Album = meta.Album.Name
			if len(meta.Album.Images) > 0 {
				job.ArtURL = meta.Album.Images[0].URL
			}
		}
		jobs = append(jobs, job)
	}

	total := len(jobs)
	ui.LogAlways("\n--- Starting downloads (%d total) ---\n\n", total)

	jobsChan := make(chan DownloadJob, total)
	resultsChan := make(chan Result, total)

	for w := 1; w <= 3; w++ {
		go func() {
			for job := range jobsChan {
				err := downloader.DownloadTrack(job.Artist, job.Title, outputDir, job.Album, job.ArtURL)
				resultsChan <- Result{Job: job, Err: err}
			}
		}()
	}

	for _, j := range jobs { jobsChan <- j }
	close(jobsChan)

	success, skipped, failed := 0, 0, []string{}
	for i := 0; i < total; i++ {
		res := <-resultsChan
		if res.Err != nil {
			if errors.Is(res.Err, downloader.ErrSkipped) {
				skipped++
				ui.LogInfo("[-] Skipped: %s - %s\n", res.Job.Artist, res.Job.Title)
			} else {
				failed = append(failed, fmt.Sprintf("%s - %s", res.Job.Artist, res.Job.Title))
				ui.LogInfo("[X] Failed:  %s - %s\n", res.Job.Artist, res.Job.Title)
			}
		} else {
			success++
			ui.LogInfo("[+] %s - %s\n", res.Job.Artist, res.Job.Title)
		}
	}

	ui.LogAlways("\n--- Summary: %d Downloaded, %d Skipped, %d Failed ---\n", success, skipped, len(failed))
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().IntVarP(&count, "count", "c", 10, "Number of similar tracks")
	rootCmd.Flags().StringVarP(&outputDir, "output", "o", "./foraged-tracks", "Output directory")
	rootCmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "Minimal output")
	rootCmd.Flags().BoolVar(&only, "only", false, "Only download the provided track")
	rootCmd.Flags().BoolVar(&includeSource, "include-source", false, "Include source track")
	rootCmd.Flags().StringVar(&textInput, "text", "", "Search via text")
}