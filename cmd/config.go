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

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage AWS profiles",
}

type ProfileConfig struct {
	Name            string
	Region          string
	AuthType        string
	SSOStartURL     string
	SSORegion       string
	SSOAccountID    string
	SSORoleName     string
	AccessKeyID     string
	SecretAccessKey string
}

var regions = []string{
	"us-east-1", "us-east-2", "us-west-1", "us-west-2",
	"eu-west-1", "eu-west-2", "eu-central-1",
	"ap-southeast-1", "ap-southeast-2", "ap-northeast-1",
}

var configLsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List all AWS profiles",
	Run: func(cmd *cobra.Command, args []string) {
		configPath := filepath.Join(os.Getenv("HOME"), ".aws", "config")
		credPath := filepath.Join(os.Getenv("HOME"), ".aws", "credentials")

		profiles := make(map[string]bool)

		if cfg, err := ini.Load(configPath); err == nil {
			for _, section := range cfg.Sections() {
				name := section.Name()
				if name == ini.DefaultSection {
					profiles["default"] = true
				} else if strings.HasPrefix(name, "profile ") {
					profiles[strings.TrimPrefix(name, "profile ")] = true
				}
			}
		}

		if creds, err := ini.Load(credPath); err == nil {
			for _, section := range creds.Sections() {
				name := section.Name()
				if name == ini.DefaultSection {
					profiles["default"] = true
				} else {
					profiles[name] = true
				}
			}
		}

		fmt.Println("✨ Available AWS profiles:")
		for name := range profiles {
			fmt.Printf("   %s\n", name)
		}
	},
}

var configAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new AWS profile",
	Run: func(cmd *cobra.Command, args []string) {
		profile := &ProfileConfig{}

		// Profile name
		namePrompt := &survey.Input{
			Message: "Profile name:",
			Help:    "Enter a unique name for this profile",
		}
		survey.AskOne(namePrompt, &profile.Name, survey.WithValidator(survey.Required))

		// Region selection
		regionPrompt := &survey.Select{
			Message: "Select AWS region:",
			Options: regions,
			Default: "us-east-1",
		}
		survey.AskOne(regionPrompt, &profile.Region)

		// Authentication type
		authTypePrompt := &survey.Select{
			Message: "Choose authentication method:",
			Options: []string{"AWS SSO", "Access Keys"},
			Default: "AWS SSO",
		}
		survey.AskOne(authTypePrompt, &profile.AuthType)

		if profile.AuthType == "AWS SSO" {
			ssoQuestions := []*survey.Question{
				{
					Name: "SSOStartURL",
					Prompt: &survey.Input{
						Message: "SSO start URL:",
						Help:    "Enter your AWS SSO start URL",
					},
					Validate: survey.Required,
				},
				{
					Name: "SSORegion",
					Prompt: &survey.Select{
						Message: "SSO region:",
						Options: regions,
						Default: profile.Region,
					},
				},
				{
					Name: "SSOAccountID",
					Prompt: &survey.Input{
						Message: "AWS Account ID:",
						Help:    "Enter your 12-digit AWS account ID",
					},
					Validate: func(val interface{}) error {
						str, _ := val.(string)
						if len(str) != 12 {
							return fmt.Errorf("account ID must be 12 digits")
						}
						return nil
					},
				},
				{
					Name: "SSORoleName",
					Prompt: &survey.Input{
						Message: "SSO Role name:",
						Help:    "Enter the IAM role name for SSO",
					},
					Validate: survey.Required,
				},
			}
			survey.Ask(ssoQuestions, profile)
		} else {
			accessKeyQuestions := []*survey.Question{
				{
					Name: "AccessKeyID",
					Prompt: &survey.Input{
						Message: "AWS Access Key ID:",
					},
					Validate: survey.Required,
				},
				{
					Name: "SecretAccessKey",
					Prompt: &survey.Password{
						Message: "AWS Secret Access Key:",
					},
					Validate: survey.Required,
				},
			}
			survey.Ask(accessKeyQuestions, profile)
		}

		// Save the profile
		configPath := filepath.Join(os.Getenv("HOME"), ".aws", "config")
		credPath := filepath.Join(os.Getenv("HOME"), ".aws", "credentials")

		cfg, _ := ini.LooseLoad(configPath)
		creds, _ := ini.LooseLoad(credPath)

		sectionName := "profile " + profile.Name
		section := cfg.Section(sectionName)
		section.Key("region").SetValue(profile.Region)

		if profile.AuthType == "AWS SSO" {
			section.Key("sso_start_url").SetValue(profile.SSOStartURL)
			section.Key("sso_region").SetValue(profile.SSORegion)
			section.Key("sso_account_id").SetValue(profile.SSOAccountID)
			section.Key("sso_role_name").SetValue(profile.SSORoleName)
		} else {
			credSection := creds.Section(profile.Name)
			credSection.Key("aws_access_key_id").SetValue(profile.AccessKeyID)
			credSection.Key("aws_secret_access_key").SetValue(profile.SecretAccessKey)
		}

		if err := cfg.SaveTo(configPath); err != nil {
			log.Fatalf("Failed to write config file: %v", err)
		}
		if err := creds.SaveTo(credPath); err != nil {
			log.Fatalf("Failed to write credentials file: %v", err)
		}

		fmt.Printf("✨ Profile '%s' created successfully!\n", profile.Name)
	},
}

var configRemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove an AWS profile",
	Run: func(cmd *cobra.Command, args []string) {
		configPath := filepath.Join(os.Getenv("HOME"), ".aws", "config")
		credPath := filepath.Join(os.Getenv("HOME"), ".aws", "credentials")

		// Get existing profiles
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
			fmt.Println("❌ No profiles found to remove")
			return
		}

		var selectedProfile string
		prompt := &survey.Select{
			Message: "Choose a profile to remove:",
			Options: profiles,
		}
		survey.AskOne(prompt, &selectedProfile)

		var confirm bool
		confirmPrompt := &survey.Confirm{
			Message: fmt.Sprintf("⚠️  Are you sure you want to remove profile '%s'?", selectedProfile),
			Default: false,
		}

		survey.AskOne(confirmPrompt, &confirm)
		if !confirm {
			fmt.Println("Operation cancelled")
			return
		}

		cfg, _ := ini.LooseLoad(configPath)
		creds, _ := ini.LooseLoad(credPath)

		sectionName := "profile " + selectedProfile
		if cfg.HasSection(sectionName) {
			cfg.DeleteSection(sectionName)
		}
		if creds.HasSection(selectedProfile) {
			creds.DeleteSection(selectedProfile)
		}

		if err := cfg.SaveTo(configPath); err != nil {
			log.Fatalf("Failed to write config file: %v", err)
		}
		if err := creds.SaveTo(credPath); err != nil {
			log.Fatalf("Failed to write credentials file: %v", err)
		}

		fmt.Printf("✨ Profile '%s' has been removed\n", selectedProfile)
	},
}

var configEditCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit an existing AWS profile",
	Run: func(cmd *cobra.Command, args []string) {
		configPath := filepath.Join(os.Getenv("HOME"), ".aws", "config")
		credPath := filepath.Join(os.Getenv("HOME"), ".aws", "credentials")

		// Get existing profiles
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
			fmt.Println("❌ No profiles found to edit")
			return
		}

		var selectedProfile string
		prompt := &survey.Select{
			Message: "Choose a profile to edit:",
			Options: profiles,
		}
		survey.AskOne(prompt, &selectedProfile)

		// Load existing configuration
		cfg, _ := ini.LooseLoad(configPath)
		creds, _ := ini.LooseLoad(credPath)
		sectionName := "profile " + selectedProfile
		section := cfg.Section(sectionName)

		// Determine what to edit
		editOptions := []string{"Region"}
		if section.HasKey("sso_start_url") {
			editOptions = append(editOptions, "SSO Configuration")
		}
		if creds.HasSection(selectedProfile) {
			editOptions = append(editOptions, "Access Keys")
		}

		var editChoice string
		editPrompt := &survey.Select{
			Message: "What would you like to edit?",
			Options: editOptions,
		}
		survey.AskOne(editPrompt, &editChoice)

		switch editChoice {
		case "Region":
			var newRegion string
			regionPrompt := &survey.Select{
				Message: "Select new AWS region:",
				Options: regions,
				Default: section.Key("region").Value(),
			}
			survey.AskOne(regionPrompt, &newRegion)
			section.Key("region").SetValue(newRegion)

		case "SSO Configuration":
			ssoQuestions := []*survey.Question{
				{
					Name: "ssoStartURL",
					Prompt: &survey.Input{
						Message: "New SSO start URL:",
						Default: section.Key("sso_start_url").Value(),
					},
				},
				{
					Name: "ssoRegion",
					Prompt: &survey.Select{
						Message: "New SSO region:",
						Options: regions,
						Default: section.Key("sso_region").Value(),
					},
				},
				{
					Name: "ssoAccountID",
					Prompt: &survey.Input{
						Message: "New AWS Account ID:",
						Default: section.Key("sso_account_id").Value(),
					},
					Validate: func(val interface{}) error {
						str, _ := val.(string)
						if len(str) != 12 {
							return fmt.Errorf("account ID must be 12 digits")
						}
						return nil
					},
				},
				{
					Name: "ssoRoleName",
					Prompt: &survey.Input{
						Message: "New SSO Role name:",
						Default: section.Key("sso_role_name").Value(),
					},
				},
			}

			answers := struct {
				SSOStartURL  string
				SSORegion    string
				SSOAccountID string
				SSORoleName  string
			}{}

			survey.Ask(ssoQuestions, &answers)
			section.Key("sso_start_url").SetValue(answers.SSOStartURL)
			section.Key("sso_region").SetValue(answers.SSORegion)
			section.Key("sso_account_id").SetValue(answers.SSOAccountID)
			section.Key("sso_role_name").SetValue(answers.SSORoleName)

		case "Access Keys":
			credSection := creds.Section(selectedProfile)
			accessKeyQuestions := []*survey.Question{
				{
					Name: "accessKeyID",
					Prompt: &survey.Input{
						Message: "New AWS Access Key ID:",
						Default: credSection.Key("aws_access_key_id").Value(),
					},
				},
				{
					Name: "secretAccessKey",
					Prompt: &survey.Password{
						Message: "New AWS Secret Access Key:",
					},
				},
			}

			answers := struct {
				AccessKeyID     string
				SecretAccessKey string
			}{}

			survey.Ask(accessKeyQuestions, &answers)
			credSection.Key("aws_access_key_id").SetValue(answers.AccessKeyID)
			if answers.SecretAccessKey != "" {
				credSection.Key("aws_secret_access_key").SetValue(answers.SecretAccessKey)
			}
		}

		if err := cfg.SaveTo(configPath); err != nil {
			log.Fatalf("Failed to write config file: %v", err)
		}
		if err := creds.SaveTo(credPath); err != nil {
			log.Fatalf("Failed to write credentials file: %v", err)
		}

		fmt.Printf("✨ Profile '%s' updated successfully!\n", selectedProfile)
	},
}

func init() {
	configCmd.AddCommand(configLsCmd)
	configCmd.AddCommand(configAddCmd)
	configCmd.AddCommand(configRemoveCmd)
	configCmd.AddCommand(configEditCmd)
	rootCmd.AddCommand(configCmd)
}
