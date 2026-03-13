package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
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
}

func GetConfigPath() (string, error) {
	home, _ := os.UserHomeDir()
	configDir := filepath.Join(home, ".config", "forage")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", err
	}
	return filepath.Join(configDir, "config.yaml"), nil
}

func Load() (*Config, error) {
	_ = godotenv.Load() // Load .env if it exists

	cfg := &Config{
		DefaultCount:  10,
		OutputDir:     "./foraged-tracks",
	}

	path, err := GetConfigPath()
	if err == nil {
		if data, err := os.ReadFile(path); err == nil {
			_ = yaml.Unmarshal(data, cfg)
		}
	}

	// Override with Env Vars
	if id := os.Getenv("SPOTIFY_CLIENT_ID"); id != "" { cfg.SpotifyClientID = id }
	if sec := os.Getenv("SPOTIFY_CLIENT_SECRET"); sec != "" { cfg.SpotifyClientSecret = sec }
	if key := os.Getenv("LASTFM_API_KEY"); key != "" { cfg.LastFmAPIKey = key }

	if cfg.SpotifyClientID == "" || cfg.SpotifyClientSecret == "" || cfg.LastFmAPIKey == "" {
		return nil, fmt.Errorf("missing credentials in config.yaml or .env")
	}

	return cfg, nil
}