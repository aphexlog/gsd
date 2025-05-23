package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"gopkg.in/ini.v1"
)

var switchCmd = &cobra.Command{
	Use:   "switch",
	Short: "Switch between AWS profiles interactively",
	Run: func(cmd *cobra.Command, args []string) {
		configPath := filepath.Join(os.Getenv("HOME"), ".aws", "config")
		credPath := filepath.Join(os.Getenv("HOME"), ".aws", "credentials")
		gsdProfilePath := filepath.Join(os.Getenv("HOME"), ".aws", ".gsd-current")

		// Get available profiles
		profiles := make([]string, 0)
		if cfg, err := ini.Load(configPath); err == nil {
			for _, section := range cfg.Sections() {
				name := section.Name()
				if strings.HasPrefix(name, "profile ") {
					profiles = append(profiles, strings.TrimPrefix(name, "profile "))
				}
			}
		}

		if len(profiles) == 0 {
			fmt.Println("🤖 No AWS profiles found in configuration")
			return
		}

		// Get current profile
		currentProfile := "default"
		if data, err := os.ReadFile(gsdProfilePath); err == nil {
			currentProfile = string(data)
		}

		// Create profile selection prompt
		var selectedProfile string
		prompt := &survey.Select{
			Message: "🤖 Select AWS profile:",
			Options: profiles,
			Default: currentProfile,
			Help:    "Choose the AWS profile you want to use",
		}

		// Custom styling
		surveyOpts := survey.WithIcons(func(icons *survey.IconSet) {
			icons.Question.Text = "🤖"
			icons.Question.Format = "cyan"
			icons.SelectFocus.Text = "→"
			icons.SelectFocus.Format = "cyan"
		})

		err := survey.AskOne(prompt, &selectedProfile, surveyOpts)
		if err != nil {
			if err.Error() == "interrupt" {
				fmt.Println("\n🤖 Operation cancelled")
				os.Exit(0)
			}
			log.Fatalf("🤖 Error selecting profile: %v", err)
		}

		// --- CONFIG ---
		cfg, err := ini.Load(configPath)
		if err != nil {
			log.Fatalf("🤖 Unable to load config: %v", err)
		}

		srcSectionName := "profile " + selectedProfile
		if !cfg.HasSection(srcSectionName) {
			log.Fatalf("🤖 Profile '%s' not found", selectedProfile)
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
			log.Fatalf("🤖 Unable to save config: %v", err)
		}

		// --- CREDENTIALS ---
		creds, err := ini.Load(credPath)
		if err == nil && creds.HasSection(selectedProfile) {
			srcCred := creds.Section(selectedProfile)
			dstCred := creds.Section("default")

			for _, key := range dstCred.Keys() {
				dstCred.DeleteKey(key.Name())
			}
			for _, key := range srcCred.Keys() {
				dstCred.Key(key.Name()).SetValue(key.Value())
			}

			if err := creds.SaveTo(credPath); err != nil {
				log.Fatalf("🤖 Unable to save credentials: %v", err)
			}
		}

		// --- TRACK CURRENT PROFILE ---
		if err := os.WriteFile(gsdProfilePath, []byte(selectedProfile), 0600); err != nil {
			log.Printf("🤖 Note: Could not save current profile: %v", err)
		}

		fmt.Printf("🤖 Switched to profile: '%s'\n", selectedProfile)
	},
}

func init() {
	rootCmd.AddCommand(switchCmd)
}
