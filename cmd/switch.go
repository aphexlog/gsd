package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// switchCmd represents the switch command
var switchCmd = &cobra.Command{
	Use:   "switch <profile>",
	Short: "Switches the current AWS profile by printing an export statement",
	Long: `Prints an export command to set AWS_PROFILE to the given profile.
Use it like: eval "$(gsd switch dev)"`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		profile := args[0]
		fmt.Printf("export AWS_PROFILE=%s\n", profile)
	},
}

func init() {
	rootCmd.AddCommand(switchCmd)
}
