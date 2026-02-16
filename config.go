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
			return nil, fmt.Errorf("config not found. create %s with your API credentials", configPath)
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