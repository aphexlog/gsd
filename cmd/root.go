package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gsd",
	Short: "🤖 AWS profile management assistant",
	Long: `🤖 GSD (Get Stuff Done) - Your AWS Profile Assistant
A friendly tool for managing AWS profiles and services.
Making AWS profile management simple and efficient.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Hide the completion command
	rootCmd.CompletionOptions.HiddenDefaultCmd = true
}
