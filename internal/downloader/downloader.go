package downloader

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"forage/internal/ui"
	"github.com/bogem/id3v2/v2"
)

var ErrSkipped = errors.New("skipped")

func sanitizeFilename(name string) string {
	invalid := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|"}
	result := name
	for _, char := range invalid {
		result = strings.ReplaceAll(result, char, "_")
	}
	return result
}

func DownloadTrack(artist, track, outputDir, album, albumArtURL string) error {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %v", err)
	}

	safeArtist := sanitizeFilename(artist)
	safeTrack := sanitizeFilename(track)
	filename := fmt.Sprintf("%s - %s.mp3", safeArtist, safeTrack)
	outputPath := filepath.Join(outputDir, filename)

	if _, err := os.Stat(outputPath); err == nil {
		return ErrSkipped
	}

	query := fmt.Sprintf("%s %s audio", artist, track)

	cmd := exec.Command("yt-dlp",
		"-x",
		"--audio-format", "mp3",
		"-o", outputPath,
		fmt.Sprintf("ytsearch1:%s", query))

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("download failed: %v\n%s", err, string(output))
	}

	if err := AddMetadata(outputPath, artist, track, album, albumArtURL); err != nil {
		return nil 
	}

	return nil
}

func AddMetadata(filepath, artist, title, album, albumArtURL string) error {
	tag, err := id3v2.Open(filepath, id3v2.Options{Parse: true})
	if err != nil {
		return fmt.Errorf("error opening mp3: %v", err)
	}
	defer tag.Close()

	tag.SetArtist(artist)
	tag.SetTitle(title)
	tag.SetAlbum(album)

	if albumArtURL != "" {
		resp, err := http.Get(albumArtURL)
		if err != nil {
			ui.LogInfo("⚠ Failed to download album art\n")
		} else {
			defer resp.Body.Close()
			artData, err := io.ReadAll(resp.Body)
			if err != nil {
				ui.LogInfo("⚠ Failed to read album art\n")
			} else {
				pic := id3v2.PictureFrame{
					Encoding:    id3v2.EncodingUTF8,
					MimeType:    "image/jpeg",
					PictureType: id3v2.PTFrontCover,
					Description: "Front cover",
					Picture:     artData,
				}
				tag.AddAttachedPicture(pic)
			}
		}
	}

	if err := tag.Save(); err != nil {
		return fmt.Errorf("error saving metadata: %v", err)
	}

	return nil
}