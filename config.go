package main

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	SpotifyClientID     string `yaml:"spotify_client_id"`
	SpotifyClientSecret string `yaml:"spotify_client_secret"`
	LastFmAPIKey        string `yaml:"lastfm_api_key"`
}

func getConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	
	configDir := filepath.Join(home, ".config", "forage")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", err
	}
	
	return filepath.Join(configDir, "config.yaml"), nil
}

func loadConfig() (*Config, error) {
	// Try .env first (for development)
	spotifyClientID := os.Getenv("SPOTIFY_CLIENT_ID")
	spotifyClientSecret := os.Getenv("SPOTIFY_CLIENT_SECRET")
	lastfmAPIKey := os.Getenv("LASTFM_API_KEY")
	
	if spotifyClientID != "" && spotifyClientSecret != "" && lastfmAPIKey != "" {
		return &Config{
			SpotifyClientID:     spotifyClientID,
			SpotifyClientSecret: spotifyClientSecret,
			LastFmAPIKey:        lastfmAPIKey,
		}, nil
	}
	
	// Try config file
	configPath, err := getConfigPath()
	if err != nil {
		return nil, err
	}
	
	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			if err := createConfigTemplate(); err != nil {
				return nil, fmt.Errorf("failed to create config template: %v", err)
			}
			return nil, fmt.Errorf("created config template at %s - please add your API credentials", configPath)
		}
		return nil, err
	}
	
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}
	
	if config.SpotifyClientID == "" || config.SpotifyClientSecret == "" || config.LastFmAPIKey == "" {
		return nil, fmt.Errorf("missing credentials in config file")
	}
	
	return &config, nil
}

func createConfigTemplate() error {
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}
	
	template := `# Forage configuration
# Get credentials from:
# Spotify: https://developer.spotify.com/dashboard
# Last.fm: https://www.last.fm/api/account/create

spotify_client_id: ""
spotify_client_secret: ""
lastfm_api_key: ""
`
	
	return os.WriteFile(configPath, []byte(template), 0644)
}