package cmd

import (
	"forage/internal/config"
	"forage/internal/ui"
	"os"

	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Open the configuration file",
	Run: func(cmd *cobra.Command, args []string) {
		path, err := config.GetConfigPath()
		if err != nil {
			ui.LogAlways("Error: %v\n", err)
			os.Exit(1)
		}

		if _, err := os.Stat(path); os.IsNotExist(err) {
			template := `
spotify_client_id: ""
spotify_client_secret: ""
lastfm_api_key: ""
default_count: 10
output_dir: "./foraged-tracks"
quiet_mode: false
include_source: false
`
			_ = os.WriteFile(path, []byte(template), 0644)
			ui.LogAlways("✓ Created template at: %s\n", path)
		}

		ui.LogAlways("Opening config...\n")
		ui.OpenFile(path)
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}