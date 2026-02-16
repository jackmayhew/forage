package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type SpotifyAuthResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

type SpotifyTrack struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Artists []struct {
		Name string `json:"name"`
	} `json:"artists"`
	Album struct {
		Name   string `json:"name"`
		Images []struct {
			URL string `json:"url"`
		} `json:"images"`
	} `json:"album"`
}

func getSpotifyToken(clientID, clientSecret string) (string, error) {
	authURL := "https://accounts.spotify.com/api/token"

	data := url.Values{}
	data.Set("grant_type", "client_credentials")

	req, err := http.NewRequest("POST", authURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(clientID, clientSecret)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var authResp SpotifyAuthResponse
	err = json.Unmarshal(body, &authResp)
	if err != nil {
		return "", err
	}

	return authResp.AccessToken, nil
}

func getTrackInfo(token, trackID string) (*SpotifyTrack, error) {
	trackURL := fmt.Sprintf("https://api.spotify.com/v1/tracks/%s", trackID)

	req, err := http.NewRequest("GET", trackURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var track SpotifyTrack
	err = json.Unmarshal(body, &track)
	if err != nil {
		return nil, err
	}

	return &track, nil
}

func getTrackInfoBySearch(token, artist, track string) (*SpotifyTrack, error) {
	query := fmt.Sprintf("track:%s artist:%s", track, artist)
	searchURL := fmt.Sprintf("https://api.spotify.com/v1/search?q=%s&type=track&limit=1", 
		url.QueryEscape(query))

	req, err := http.NewRequest("GET", searchURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result struct {
		Tracks struct {
			Items []SpotifyTrack `json:"items"`
		} `json:"tracks"`
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	if len(result.Tracks.Items) == 0 {
		return nil, fmt.Errorf("track not found")
	}

	return &result.Tracks.Items[0], nil
}

func searchSpotifyTrack(token, query string) (*SpotifyTrack, error) {
	searchURL := fmt.Sprintf("https://api.spotify.com/v1/search?q=%s&type=track&limit=1", url.QueryEscape(query))

	req, err := http.NewRequest("GET", searchURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("spotify search API responded with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var result struct {
		Tracks struct {
			Items []SpotifyTrack `json:"items"`
		} `json:"tracks"`
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	if len(result.Tracks.Items) == 0 {
		return nil, fmt.Errorf("no track found for '%s'", query)
	}

	return &result.Tracks.Items[0], nil
}