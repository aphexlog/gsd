package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Test AWS config loading",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadDefaultConfig(context.TODO())
		if err != nil {
			log.Fatalf("unable to load AWS SDK config: %v", err)
		}

		fmt.Println("AWS config loaded successfully.")
		fmt.Println("Region:", cfg.Region)
		fmt.Println("Credentials:", cfg.Credentials)
		creds, err := cfg.Credentials.Retrieve(context.TODO())
		if err != nil {
			log.Fatalf("unable to retrieve credentials: %v", err)
		}
		println("Access Key ID:", creds.AccessKeyID)
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
