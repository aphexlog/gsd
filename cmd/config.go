package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
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
		creds, err := cfg.Credentials.Retrieve(context.TODO())
		if err != nil {
			log.Fatalf("unable to retrieve credentials: %v", err)
		}
		fmt.Println("Access Key ID:", creds.AccessKeyID)

		// Add a profile called test to the config
		customCfg := aws.Config{
			Region: "us-east-1",
			Credentials: credentials.NewStaticCredentialsProvider(
				"test-access-key-id",
				"test-secret-access-key",
				""),
		}
		fmt.Println("Custom AWS config created successfully.")
		fmt.Println("Custom Region:", customCfg.Region)
		customCreds, err := customCfg.Credentials.Retrieve(context.TODO())
		if err != nil {
			log.Fatalf("unable to retrieve custom credentials: %v", err)
		}
		fmt.Println("Custom Access Key ID:", customCreds.AccessKeyID)
	},
}

func init() {
	rootCmd.AddCommand(configCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// configCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
