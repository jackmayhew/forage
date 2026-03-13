package lastfm

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type SimilarResponse struct {
	SimilarTracks struct {
		Track []Track `json:"track"`
	} `json:"similartracks"`
}

type Track struct {
	Name   string `json:"name"`
	Artist struct {
		Name string `json:"name"`
	} `json:"artist"`
}

func GetSimilarTracks(apiKey, artist, track string, limit int) ([]Track, error) {
	baseURL := "http://ws.audioscrobbler.com/2.0/"
	params := url.Values{}
	params.Set("method", "track.getsimilar")
	params.Set("artist", artist)
	params.Set("track", track)
	params.Set("api_key", apiKey)
	params.Set("limit", fmt.Sprintf("%d", limit))
	params.Set("format", "json")

	fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	resp, err := http.Get(fullURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result SimilarResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return result.SimilarTracks.Track, nil
}