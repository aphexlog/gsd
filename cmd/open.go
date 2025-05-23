package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"gopkg.in/ini.v1"
)

// AWS services and their console URLs
var awsServices = map[string]string{
	"Console (Main)":    "https://signin.aws.amazon.com/console",
	"SSO":              "https://d-{account_id}.awsapps.com/start",  // Will be replaced with actual SSO URL
	"EC2":              "https://console.aws.amazon.com/ec2/v2/home",
	"S3":               "https://s3.console.aws.amazon.com/s3/home",
	"Lambda":           "https://console.aws.amazon.com/lambda/home",
	"CloudFormation":   "https://console.aws.amazon.com/cloudformation/home",
	"CloudWatch":       "https://console.aws.amazon.com/cloudwatch/home",
	"IAM":             "https://console.aws.amazon.com/iam/home",
	"RDS":             "https://console.aws.amazon.com/rds/home",
	"DynamoDB":        "https://console.aws.amazon.com/dynamodb/home",
	"ECS":             "https://console.aws.amazon.com/ecs/home",
	"EKS":             "https://console.aws.amazon.com/eks/home",
	"API Gateway":      "https://console.aws.amazon.com/apigateway/home",
	"Route 53":        "https://console.aws.amazon.com/route53/home",
	"SQS":             "https://console.aws.amazon.com/sqs/home",
	"SNS":             "https://console.aws.amazon.com/sns/home",
	"Secrets Manager": "https://console.aws.amazon.com/secretsmanager/home",
	"Systems Manager": "https://console.aws.amazon.com/systems-manager/home",
	"CodePipeline":    "https://console.aws.amazon.com/codesuite/codepipeline/home",
	"CodeBuild":       "https://console.aws.amazon.com/codesuite/codebuild/home",
	"Amplify":         "https://console.aws.amazon.com/amplify/home",
}

// openCmd represents the open command
var openCmd = &cobra.Command{
	Use:   "open",
	Short: "Open AWS Console or specific services",
	Long:  "🤖 Select and open AWS Console or services in your browser",
	Run: func(cmd *cobra.Command, args []string) {
		// Get available profiles
		configPath := filepath.Join(os.Getenv("HOME"), ".aws", "config")
		profiles := []string{"default"}

		if cfg, err := ini.Load(configPath); err == nil {
			for _, section := range cfg.Sections() {
				name := section.Name()
				if strings.HasPrefix(name, "profile ") {
					profiles = append(profiles, strings.TrimPrefix(name, "profile "))
				}
			}
		}

		// Get current profile
		currentProfile := "default"
		gsdProfilePath := filepath.Join(os.Getenv("HOME"), ".aws", ".gsd-current")
		if data, err := os.ReadFile(gsdProfilePath); err == nil {
			currentProfile = string(data)
		}

		// Create questions
		questions := []*survey.Question{
			{
				Name: "service",
				Prompt: &survey.Select{
					Message: "🤖 Select AWS service:",
					Options: getServicesList(),
					Default: "Console (Main)",
				},
			},
			{
				Name: "profile",
				Prompt: &survey.Select{
					Message: "🤖 Select AWS profile:",
					Options: profiles,
					Default: currentProfile,
				},
			},
		}

		// Custom styling
		surveyOpts := survey.WithIcons(func(icons *survey.IconSet) {
			icons.Question.Text = "🤖"
			icons.Question.Format = "cyan"
			icons.SelectFocus.Text = "→"
			icons.SelectFocus.Format = "cyan"
		})

		// Get answers
		answers := struct {
			Service string
			Profile string
		}{}

		err := survey.Ask(questions, &answers, surveyOpts)
		if err != nil {
			if err.Error() == "interrupt" {
				fmt.Println("\n🤖 Operation cancelled")
				os.Exit(0)
			}
			log.Fatalf("🤖 Error getting input: %v", err)
		}

		// Get the URL for the selected service
		url := awsServices[answers.Service]
		if url == "" {
			log.Fatalf("🤖 Service not found")
		}

		// If SSO is selected, get the SSO URL from the profile configuration
		if answers.Service == "SSO" {
			cfg, err := ini.Load(configPath)
			if err != nil {
				log.Fatalf("🤖 Unable to load AWS config: %v", err)
			}

			sectionName := "profile " + answers.Profile
			if answers.Profile == "default" {
				sectionName = "default"
			}

			section := cfg.Section(sectionName)
			
			// First try to get the sso_session name
			ssoSession := section.Key("sso_session").String()
			var ssoStartURL string
			
			if ssoSession != "" {
				// Look up the SSO session configuration
				ssoSectionName := "sso-session " + ssoSession
				if ssoSection := cfg.Section(ssoSectionName); ssoSection != nil {
					ssoStartURL = ssoSection.Key("sso_start_url").String()
				}
			} else {
				// Fallback to looking for sso_start_url directly in the profile (legacy configuration)
				ssoStartURL = section.Key("sso_start_url").String()
			}
			
			if ssoStartURL == "" {
				log.Fatalf("🤖 No SSO configuration found for profile '%s'", answers.Profile)
			}
			
			url = ssoStartURL
		}

		// Open the URL in the default browser
		err = openBrowser(url)
		if err != nil {
			log.Fatalf("🤖 Unable to open browser: %v", err)
		}

		fmt.Printf("🤖 Opening %s for profile '%s'\n", answers.Service, answers.Profile)
	},
}

func init() {
	rootCmd.AddCommand(openCmd)
}

// getServicesList returns a sorted list of AWS service names
func getServicesList() []string {
	services := make([]string, 0, len(awsServices))
	for service := range awsServices {
		services = append(services, service)
	}
	return services
}

// openBrowser opens the specified URL in the default browser
func openBrowser(url string) error {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	return err
}
