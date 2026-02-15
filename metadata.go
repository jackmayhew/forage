package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/bogem/id3v2/v2"
)

func addMetadata(filepath, artist, title, album, albumArtURL string) error {
	tag, err := id3v2.Open(filepath, id3v2.Options{Parse: true})
	if err != nil {
		return fmt.Errorf("error opening mp3: %v", err)
	}
	defer tag.Close()

	tag.SetArtist(artist)
	tag.SetTitle(title)
	tag.SetAlbum(album)

	// Download and add album art
	if albumArtURL != "" {
		resp, err := http.Get(albumArtURL)
		if err == nil {
			defer resp.Body.Close()
			artData, err := io.ReadAll(resp.Body)
			if err == nil {
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