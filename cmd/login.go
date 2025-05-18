package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Logs in to AWS using the current profile",
	Long:  `Logs in to AWS using the profile set in the AWS_PROFILE environment variable.`,
	Run: func(cmd *cobra.Command, args []string) {
		profile := os.Getenv("AWS_PROFILE")
		if profile == "" {
			profile = "default"
		}

		_, err := exec.LookPath("aws")
		if err != nil {
			fmt.Println("❌ AWS CLI is not installed or not found in PATH.")
			os.Exit(1)
		}

		fmt.Printf("🔐 Logging in with profile '%s'...\n", profile)

		cmdExec := exec.Command("aws", "sso", "login", "--profile", profile)
		cmdExec.Stdout = os.Stdout
		cmdExec.Stderr = os.Stderr
		cmdExec.Stdin = os.Stdin

		if err := cmdExec.Run(); err != nil {
			fmt.Printf("❌ Login failed: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("✅ Login successful.")
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
