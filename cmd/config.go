package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/ini.v1"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage AWS profiles",
	// Additional configuration...
}

var configLsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List all AWS profiles",
	Run: func(cmd *cobra.Command, args []string) {
		configPath := filepath.Join(os.Getenv("HOME"), ".aws", "config")
		credPath := filepath.Join(os.Getenv("HOME"), ".aws", "credentials")

		profiles := make(map[string]bool)

		// Parse config
		if cfg, err := ini.Load(configPath); err == nil {
			for _, section := range cfg.Sections() {
				name := section.Name()
				if name == ini.DefaultSection {
					profiles["default"] = true
				} else if strings.HasPrefix(name, "profile ") {
					profiles[strings.TrimPrefix(name, "profile ")] = true
				}
			}
		} else {
			log.Printf("Warning: could not load config file: %v", err)
		}

		// Parse credentials
		if creds, err := ini.Load(credPath); err == nil {
			for _, section := range creds.Sections() {
				name := section.Name()
				if name == ini.DefaultSection {
					profiles["default"] = true
				} else {
					profiles[name] = true
				}
			}
		} else {
			log.Printf("Warning: could not load credentials file: %v", err)
		}

		// Print results
		fmt.Println("Available AWS profiles:")
		for name := range profiles {
			fmt.Printf(" - %s\n", name)
		}
	},
}

func init() {
	configCmd.AddCommand(configLsCmd)
	rootCmd.AddCommand(configCmd)
}
