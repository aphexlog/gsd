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
		} else {
			log.Printf("Warning: could not load config file: %v", err)
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
		} else {
			log.Printf("Warning: could not load credentials file: %v", err)
		}

		fmt.Println("Available AWS profiles:")
		for name := range profiles {
			fmt.Printf(" - %s\n", name)
		}
	},
}

var configAddCmd = &cobra.Command{
	Use:   "add <profile-name>",
	Short: "Add a new AWS profile",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		profile := args[0]
		region, _ := cmd.Flags().GetString("region")
		ssoStartURL, _ := cmd.Flags().GetString("sso-start-url")
		ssoRegion, _ := cmd.Flags().GetString("sso-region")
		ssoAccountID, _ := cmd.Flags().GetString("sso-account-id")
		ssoRoleName, _ := cmd.Flags().GetString("sso-role-name")
		accessKeyID, _ := cmd.Flags().GetString("access-key-id")
		secretAccessKey, _ := cmd.Flags().GetString("secret-access-key")

		configPath := filepath.Join(os.Getenv("HOME"), ".aws", "config")
		credPath := filepath.Join(os.Getenv("HOME"), ".aws", "credentials")

		cfg, _ := ini.LooseLoad(configPath)
		creds, _ := ini.LooseLoad(credPath)

		sectionName := "profile " + profile
		section := cfg.Section(sectionName)

		if region != "" {
			section.Key("region").SetValue(region)
		}

		if ssoStartURL != "" {
			// Validate required SSO fields
			if ssoRegion == "" || ssoAccountID == "" || ssoRoleName == "" {
				log.Fatalf("❌ Missing required SSO fields: sso-region, sso-account-id, and sso-role-name must be provided with sso-start-url")
			}
			section.Key("sso_start_url").SetValue(ssoStartURL)
			section.Key("sso_region").SetValue(ssoRegion)
			section.Key("sso_account_id").SetValue(ssoAccountID)
			section.Key("sso_role_name").SetValue(ssoRoleName)
			section.Key("credential_process").SetValue("")
		}

		if accessKeyID != "" && secretAccessKey != "" {
			credSection := creds.Section(profile)
			credSection.Key("aws_access_key_id").SetValue(accessKeyID)
			credSection.Key("aws_secret_access_key").SetValue(secretAccessKey)
		}

		if err := cfg.SaveTo(configPath); err != nil {
			log.Fatalf("Failed to write config file: %v", err)
		}
		if err := creds.SaveTo(credPath); err != nil {
			log.Fatalf("Failed to write credentials file: %v", err)
		}

		fmt.Printf("✅ Profile '%s' added successfully.\n", profile)
	},
}

var configRemoveCmd = &cobra.Command{
	Use:   "remove <profile-name>",
	Short: "Remove an AWS profile from config and credentials",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		profile := args[0]
		configPath := filepath.Join(os.Getenv("HOME"), ".aws", "config")
		credPath := filepath.Join(os.Getenv("HOME"), ".aws", "credentials")

		fmt.Printf("⚠️  You are about to erase the AWS profile '%s' from existence.\n", profile)
		fmt.Print("Are you *absolutely* sure? This action cannot be undone (y/N): ")

		var input string
		fmt.Scanln(&input)
		if strings.ToLower(input) != "y" && strings.ToLower(input) != "yes" {
			fmt.Println("🧼 Crisis averted. Profile is safe... for now.")
			return
		}

		cfg, _ := ini.LooseLoad(configPath)
		creds, _ := ini.LooseLoad(credPath)

		sectionName := "profile " + profile
		if cfg.HasSection(sectionName) {
			cfg.DeleteSection(sectionName)
			fmt.Printf("💀 Deleted profile '%s' from config.\n", profile)
		}
		if creds.HasSection(profile) {
			creds.DeleteSection(profile)
			fmt.Printf("🩸 Deleted profile '%s' from credentials.\n", profile)
		}

		if err := cfg.SaveTo(configPath); err != nil {
			log.Fatalf("Failed to write config file: %v", err)
		}
		if err := creds.SaveTo(credPath); err != nil {
			log.Fatalf("Failed to write credentials file: %v", err)
		}

		fmt.Printf("☠️  It's done. '%s' has been removed.\n", profile)
	},
}

var configEditCmd = &cobra.Command{
	Use:   "edit <profile-name>",
	Short: "Edit an existing AWS profile",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		profile := args[0]
		configPath := filepath.Join(os.Getenv("HOME"), ".aws", "config")
		credPath := filepath.Join(os.Getenv("HOME"), ".aws", "credentials")

		region, _ := cmd.Flags().GetString("region")
		ssoStartURL, _ := cmd.Flags().GetString("sso-start-url")
		ssoRegion, _ := cmd.Flags().GetString("sso-region")
		ssoAccountID, _ := cmd.Flags().GetString("sso-account-id")
		ssoRoleName, _ := cmd.Flags().GetString("sso-role-name")
		accessKeyID, _ := cmd.Flags().GetString("access-key-id")
		secretAccessKey, _ := cmd.Flags().GetString("secret-access-key")

		cfg, _ := ini.LooseLoad(configPath)
		creds, _ := ini.LooseLoad(credPath)

		sectionName := "profile " + profile
		if !cfg.HasSection(sectionName) && !creds.HasSection(profile) {
			fmt.Printf("❌ Profile '%s' not found in config or credentials.\n", profile)
			return
		}

		section := cfg.Section(sectionName)
		if region != "" {
			section.Key("region").SetValue(region)
		}
		if ssoStartURL != "" {
			if ssoRegion == "" || ssoAccountID == "" || ssoRoleName == "" {
				log.Fatalf("❌ Missing required SSO fields: sso-region, sso-account-id, and sso-role-name must be provided with sso-start-url")
			}
			section.Key("sso_start_url").SetValue(ssoStartURL)
			section.Key("sso_region").SetValue(ssoRegion)
			section.Key("sso_account_id").SetValue(ssoAccountID)
			section.Key("sso_role_name").SetValue(ssoRoleName)
			section.Key("credential_process").SetValue("")
		}

		if accessKeyID != "" || secretAccessKey != "" {
			credSection := creds.Section(profile)
			if accessKeyID != "" {
				credSection.Key("aws_access_key_id").SetValue(accessKeyID)
			}
			if secretAccessKey != "" {
				credSection.Key("aws_secret_access_key").SetValue(secretAccessKey)
			}
		}

		if err := cfg.SaveTo(configPath); err != nil {
			log.Fatalf("Failed to write config file: %v", err)
		}
		if err := creds.SaveTo(credPath); err != nil {
			log.Fatalf("Failed to write credentials file: %v", err)
		}

		fmt.Printf("🛠️  Profile '%s' updated.\n", profile)
	},
}

func init() {
	configCmd.AddCommand(configLsCmd)
	configCmd.AddCommand(configAddCmd)
	configCmd.AddCommand(configRemoveCmd)
	configCmd.AddCommand(configEditCmd)
	rootCmd.AddCommand(configCmd)

	// Shared flags
	configAddCmd.Flags().String("region", "", "AWS region for the profile")
	configAddCmd.Flags().String("sso-start-url", "", "SSO start URL")
	configAddCmd.Flags().String("sso-region", "", "SSO region")
	configAddCmd.Flags().String("sso-account-id", "", "SSO account ID")
	configAddCmd.Flags().String("sso-role-name", "", "SSO role name")
	configAddCmd.Flags().String("access-key-id", "", "AWS access key ID")
	configAddCmd.Flags().String("secret-access-key", "", "AWS secret access key")

	configEditCmd.Flags().String("region", "", "Update AWS region")
	configEditCmd.Flags().String("sso-start-url", "", "Update SSO start URL")
	configEditCmd.Flags().String("sso-region", "", "Update SSO region")
	configEditCmd.Flags().String("sso-account-id", "", "Update SSO account ID")
	configEditCmd.Flags().String("sso-role-name", "", "Update SSO role name")
	configEditCmd.Flags().String("access-key-id", "", "Update AWS access key ID")
	configEditCmd.Flags().String("secret-access-key", "", "Update AWS secret access key")
}
