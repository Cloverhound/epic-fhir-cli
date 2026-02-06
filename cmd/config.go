package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/jbogarin/fhir-cli/internal/config"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage CLI configuration",
	Long:  `Manage fhir-cli configuration including profiles, credentials, and settings.`,
}

var configInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize configuration with interactive setup",
	RunE: func(cmd *cobra.Command, args []string) error {
		reader := bufio.NewReader(os.Stdin)

		fmt.Println("fhir-cli Configuration Setup")
		fmt.Println("=============================")
		fmt.Println()

		// Profile name
		fmt.Print("Profile name [sandbox]: ")
		profileName, _ := reader.ReadString('\n')
		profileName = strings.TrimSpace(profileName)
		if profileName == "" {
			profileName = "sandbox"
		}

		// Client ID
		fmt.Print("Client ID: ")
		clientID, _ := reader.ReadString('\n')
		clientID = strings.TrimSpace(clientID)

		// Private key source
		fmt.Println("Private key source (file path or env:VAR_NAME for environment variable)")
		fmt.Print("Private key [env:FHIR_PRIVATE_KEY]: ")
		privateKey, _ := reader.ReadString('\n')
		privateKey = strings.TrimSpace(privateKey)
		if privateKey == "" {
			privateKey = "env:FHIR_PRIVATE_KEY"
		}

		// FHIR version
		fmt.Print("FHIR version [R4]: ")
		fhirVersion, _ := reader.ReadString('\n')
		fhirVersion = strings.TrimSpace(fhirVersion)
		if fhirVersion == "" {
			fhirVersion = "R4"
		}

		// Base URL
		defaultURL := "https://fhir.epic.com/interconnect-fhir-oauth/api/FHIR/R4"
		if fhirVersion == "STU3" {
			defaultURL = "https://fhir.epic.com/interconnect-fhir-oauth/api/FHIR/STU3"
		}
		fmt.Printf("FHIR Base URL [%s]: ", defaultURL)
		baseURL, _ := reader.ReadString('\n')
		baseURL = strings.TrimSpace(baseURL)
		if baseURL == "" {
			baseURL = defaultURL
		}

		// Token URL
		defaultTokenURL := "https://fhir.epic.com/interconnect-fhir-oauth/oauth2/token"
		fmt.Printf("Token URL [%s]: ", defaultTokenURL)
		tokenURL, _ := reader.ReadString('\n')
		tokenURL = strings.TrimSpace(tokenURL)
		if tokenURL == "" {
			tokenURL = defaultTokenURL
		}

		// Scopes
		defaultScopes := "system/Patient.read system/Observation.read system/Condition.read system/MedicationRequest.read system/AllergyIntolerance.read system/Procedure.read system/Encounter.read system/DiagnosticReport.read"
		fmt.Printf("Scopes [%s]: ", defaultScopes)
		scopes, _ := reader.ReadString('\n')
		scopes = strings.TrimSpace(scopes)
		if scopes == "" {
			scopes = defaultScopes
		}

		// Create configuration
		cfg := &config.Config{
			Default: profileName,
			Profiles: map[string]config.Profile{
				profileName: {
					Name:         profileName,
					ClientID:     clientID,
					PrivateKey:   privateKey,
					TokenURL:     tokenURL,
					FHIRBaseURL:  baseURL,
					FHIRVersion:  fhirVersion,
					Scopes:       scopes,
					OutputFormat: "json",
				},
			},
		}

		if err := config.Save(cfg); err != nil {
			return fmt.Errorf("failed to save configuration: %w", err)
		}

		fmt.Printf("\nConfiguration saved to %s\n", config.GetConfigPath())
		fmt.Println("\nYou can now use fhir-cli commands. Try:")
		fmt.Println("  fhir-cli auth token    # Get an access token")
		fmt.Println("  fhir-cli patient search --name Smith")

		return nil
	},
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		if config.CurrentConfig == nil {
			return fmt.Errorf("no configuration loaded")
		}

		fmt.Printf("Configuration file: %s\n", config.GetConfigPath())
		fmt.Printf("Default profile: %s\n", config.CurrentConfig.Default)
		fmt.Printf("\nProfiles:\n")

		for name, profile := range config.CurrentConfig.Profiles {
			fmt.Printf("\n  [%s]\n", name)
			fmt.Printf("    Client ID:     %s\n", profile.ClientID)
			fmt.Printf("    Private Key:   %s\n", profile.PrivateKey)
			fmt.Printf("    FHIR Base URL: %s\n", profile.FHIRBaseURL)
			fmt.Printf("    Token URL:     %s\n", profile.TokenURL)
			fmt.Printf("    FHIR Version:  %s\n", profile.FHIRVersion)
			fmt.Printf("    Scopes:        %s\n", profile.Scopes)
		}

		return nil
	},
}

var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set a configuration value",
	Long: `Set a configuration value for the current profile.

Available keys:
  client_id      - OAuth2 client ID
  private_key    - Private key source: file path OR env:VAR_NAME
                   Examples: /path/to/key.pem, ~/keys/private.pem, env:FHIR_PRIVATE_KEY
  fhir_base_url  - FHIR API base URL
  token_url      - OAuth2 token endpoint URL
  fhir_version   - FHIR version (R4, STU3, DSTU2)
  scopes         - OAuth2 scopes
  output_format  - Default output format (json, table, yaml)`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		key := args[0]
		value := args[1]

		if config.CurrentConfig == nil || config.CurrentProfile == nil {
			return fmt.Errorf("no configuration loaded")
		}

		profileName := GetProfile()
		if profileName == "default" {
			profileName = config.CurrentConfig.Default
		}

		profile := config.CurrentConfig.Profiles[profileName]

		switch key {
		case "client_id":
			profile.ClientID = value
		case "private_key":
			// Store as-is - expansion happens at runtime in auth/jwt.go
			profile.PrivateKey = value
		case "fhir_base_url":
			profile.FHIRBaseURL = value
		case "token_url":
			profile.TokenURL = value
		case "fhir_version":
			profile.FHIRVersion = value
		case "scopes":
			profile.Scopes = value
		case "output_format":
			profile.OutputFormat = value
		default:
			return fmt.Errorf("unknown configuration key: %s", key)
		}

		config.CurrentConfig.Profiles[profileName] = profile

		if err := config.Save(config.CurrentConfig); err != nil {
			return fmt.Errorf("failed to save configuration: %w", err)
		}

		fmt.Printf("Set %s = %s for profile [%s]\n", key, value, profileName)
		return nil
	},
}

var configAddProfileCmd = &cobra.Command{
	Use:   "add-profile <name>",
	Short: "Add a new profile",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		profileName := args[0]

		if config.CurrentConfig == nil {
			config.CurrentConfig = &config.Config{
				Profiles: make(map[string]config.Profile),
			}
		}

		if _, exists := config.CurrentConfig.Profiles[profileName]; exists {
			return fmt.Errorf("profile '%s' already exists", profileName)
		}

		// Start with default profile settings
		newProfile := config.GetDefaultProfile()
		newProfile.Name = profileName

		config.CurrentConfig.Profiles[profileName] = newProfile

		if config.CurrentConfig.Default == "" {
			config.CurrentConfig.Default = profileName
		}

		if err := config.Save(config.CurrentConfig); err != nil {
			return fmt.Errorf("failed to save configuration: %w", err)
		}

		fmt.Printf("Added profile [%s]\n", profileName)
		fmt.Println("Use 'fhir-cli config set <key> <value> -p", profileName, "' to configure it")
		return nil
	},
}

var configUseCmd = &cobra.Command{
	Use:   "use <profile>",
	Short: "Set the default profile",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		profileName := args[0]

		if config.CurrentConfig == nil {
			return fmt.Errorf("no configuration loaded")
		}

		if _, exists := config.CurrentConfig.Profiles[profileName]; !exists {
			return fmt.Errorf("profile '%s' does not exist", profileName)
		}

		config.CurrentConfig.Default = profileName

		if err := config.Save(config.CurrentConfig); err != nil {
			return fmt.Errorf("failed to save configuration: %w", err)
		}

		fmt.Printf("Default profile set to [%s]\n", profileName)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configInitCmd)
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configAddProfileCmd)
	configCmd.AddCommand(configUseCmd)
}
