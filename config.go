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
	DefaultCount        int    `yaml:"default_count"`
    OutputDir           string `yaml:"output_dir"`
	QuietMode           bool   `yaml:"quiet_mode"`
    IncludeSource       bool   `yaml:"include_source"`
    UseText       		bool   `yaml:"use_text"`
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
	// defaults
	config := &Config{
		DefaultCount:  10,
		OutputDir:     "./foraged-tracks",
		QuietMode:     false,
		IncludeSource: false,
		UseText:       false,
	}

	configPath, err := getConfigPath()
	if err == nil {
		if data, err := os.ReadFile(configPath); err == nil {
			_ = yaml.Unmarshal(data, config)
		}
	}

	if envID := os.Getenv("SPOTIFY_CLIENT_ID"); envID != "" {
		config.SpotifyClientID = envID
	}
	if envSecret := os.Getenv("SPOTIFY_CLIENT_SECRET"); envSecret != "" {
		config.SpotifyClientSecret = envSecret
	}
	if envKey := os.Getenv("LASTFM_API_KEY"); envKey != "" {
		config.LastFmAPIKey = envKey
	}

	if config.SpotifyClientID == "" || config.SpotifyClientSecret == "" || config.LastFmAPIKey == "" {
		return nil, fmt.Errorf("missing credentials")
	}

	return config, nil
}

func handleConfig() {
    path, err := getConfigPath()
    if err != nil {
        logAlways("Error: %v\n", err)
        os.Exit(1)
    }

    if _, err := os.Stat(path); os.IsNotExist(err) {
        if err := createConfigTemplate(); err != nil {
            logAlways("Error creating config: %v\n", err)
            os.Exit(1)
        }
        logAlways("✓ Created config template at: %s\n", path)
    }

    logAlways("Opening config file...\n")
    openFile(path)
}

func createConfigTemplate() error {
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}
	
	template := `
# Forage configuration
# Get credentials from:
# Spotify: https://developer.spotify.com/dashboard
# Last.fm: https://www.last.fm/api/account/create

spotify_client_id: ""
spotify_client_secret: ""
lastfm_api_key: ""
default_count: 10
output_dir: "./foraged-tracks"
quiet_mode: false
include_source: false
use_text: false
`
	
	return os.WriteFile(configPath, []byte(template), 0644)
}