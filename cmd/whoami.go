package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/spf13/cobra"
)

// whoamiCmd represents the whoami command
var whoamiCmd = &cobra.Command{
	Use:   "whoami",
	Short: "Prints the current AWS profile and identity",
	Long:  `Resolves the current AWS profile and uses STS to show the active account and identity.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.TODO()

		cfg, err := config.LoadDefaultConfig(ctx)
		if err != nil {
			fmt.Printf("❌ Failed to load AWS config: %v\n", err)
			return
		}

		stsClient := sts.NewFromConfig(cfg)
		identity, err := stsClient.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
		if err != nil {
			fmt.Printf("❌ Failed to get identity: %v\n", err)
			return
		}

		profile := os.Getenv("AWS_PROFILE")
		if profile == "" {
			profile = "default"
		}

		fmt.Printf("🧠 Profile: %s\n", profile)
		fmt.Printf("🪪 Account: %s\n", *identity.Account)
		fmt.Printf("👤 ARN:     %s\n", *identity.Arn)
		fmt.Printf("🆔 User ID: %s\n", *identity.UserId)
	},
}

func init() {
	rootCmd.AddCommand(whoamiCmd)
}
