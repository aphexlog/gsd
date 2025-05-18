package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"gopkg.in/ini.v1"
)

var switchCmd = &cobra.Command{
	Use:   "switch <profile>",
	Short: "Switches the default AWS profile by rewriting the default config/credentials",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		profile := args[0]
		configPath := filepath.Join(os.Getenv("HOME"), ".aws", "config")
		credPath := filepath.Join(os.Getenv("HOME"), ".aws", "credentials")
		gsdProfilePath := filepath.Join(os.Getenv("HOME"), ".aws", ".gsd-current")

		// --- CONFIG ---
		cfg, err := ini.Load(configPath)
		if err != nil {
			log.Fatalf("❌ Failed to load config: %v", err)
		}

		srcSectionName := "profile " + profile
		if !cfg.HasSection(srcSectionName) {
			log.Fatalf("❌ Profile '%s' not found in config", profile)
		}

		src := cfg.Section(srcSectionName)
		dst := cfg.Section("default")

		for _, key := range dst.Keys() {
			dst.DeleteKey(key.Name())
		}
		for _, key := range src.Keys() {
			dst.Key(key.Name()).SetValue(key.Value())
		}

		if err := cfg.SaveTo(configPath); err != nil {
			log.Fatalf("❌ Failed to write config: %v", err)
		}

		// --- CREDENTIALS ---
		creds, err := ini.Load(credPath)
		if err == nil && creds.HasSection(profile) {
			srcCred := creds.Section(profile)
			dstCred := creds.Section("default")

			for _, key := range dstCred.Keys() {
				dstCred.DeleteKey(key.Name())
			}
			for _, key := range srcCred.Keys() {
				dstCred.Key(key.Name()).SetValue(key.Value())
			}

			if err := creds.SaveTo(credPath); err != nil {
				log.Fatalf("❌ Failed to write credentials: %v", err)
			}
		}

		// --- TRACK CURRENT PROFILE ---
		if err := os.WriteFile(gsdProfilePath, []byte(profile), 0600); err != nil {
			log.Printf("⚠️  Could not write .gsd-current file: %v", err)
		}

		fmt.Printf("✅ '%s' is now set as your default AWS profile.\n", profile)
	},
}

func init() {
	rootCmd.AddCommand(switchCmd)
}
