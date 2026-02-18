package cmd

import (
	"fmt"
	"strings"

	"github.com/jbogarin/fhir-cli/internal/auth"
	"github.com/jbogarin/fhir-cli/internal/config"
	"github.com/jbogarin/fhir-cli/internal/output"
	"github.com/spf13/cobra"
)

var (
	forceRefresh bool
	debugMode    bool
	verboseMode  bool
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authentication commands",
	Long:  `Manage OAuth2 authentication with Epic FHIR servers.`,
}

var authTokenCmd = &cobra.Command{
	Use:   "token",
	Short: "Get or refresh access token",
	Long: `Obtain an OAuth2 access token using JWT bearer assertion.

This command generates a signed JWT using your client ID and private key,
then exchanges it for an access token at the token endpoint.

The token is cached and reused until it expires. Use --force to get a new token.
Use --debug to see the JWT claims and request details for troubleshooting.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if config.CurrentProfile == nil {
			return fmt.Errorf("no profile configured. Run 'fhir-cli config init' first")
		}

		if debugMode {
			return auth.DebugAuth(config.CurrentProfile)
		}

		// Enable verbose mode to show request/response details
		auth.VerboseMode = verboseMode

		token, err := auth.GetAccessToken(config.CurrentProfile, forceRefresh)
		if err != nil {
			return fmt.Errorf("failed to get access token: %w", err)
		}

		format := output.ParseFormat(GetOutput())

		if format == output.FormatJSON {
			result := map[string]string{
				"access_token": token,
				"profile":      config.CurrentProfile.Name,
			}
			return output.Print(result, format)
		}

		fmt.Printf("Access Token: %s\n", token)
		return nil
	},
}

var authDebugCmd = &cobra.Command{
	Use:   "debug",
	Short: "Debug authentication configuration",
	Long: `Show authentication details for troubleshooting.

This displays:
- Current profile configuration
- JWT claims that would be sent
- Token endpoint URL
- Private key status (loaded/not loaded)

Common causes of unauthorized_client errors:
1. Wrong client_id - must match your Epic app registration
2. Wrong token_url - check it matches your Epic environment
3. Invalid private key - must be the key registered with Epic
4. Scopes not authorized - request only approved scopes`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if config.CurrentProfile == nil {
			return fmt.Errorf("no profile configured. Run 'fhir-cli config init' first")
		}

		return auth.DebugAuth(config.CurrentProfile)
	},
}

var authTestKeyCmd = &cobra.Command{
	Use:   "test-key",
	Short: "Test if the private key can be loaded",
	RunE: func(cmd *cobra.Command, args []string) error {
		if config.CurrentProfile == nil {
			return fmt.Errorf("no profile configured. Run 'fhir-cli config init' first")
		}

		keySource := config.CurrentProfile.PrivateKey
		fmt.Printf("Private key source: %s\n", keySource)

		if strings.HasPrefix(keySource, "env:") {
			envVar := strings.TrimPrefix(keySource, "env:")
			fmt.Printf("Environment variable: %s\n", envVar)
		}

		err := auth.TestPrivateKey(config.CurrentProfile.PrivateKey)
		if err != nil {
			fmt.Printf("❌ Failed to load private key: %v\n", err)
			return err
		}

		fmt.Println("✓ Private key loaded successfully")
		return nil
	},
}

var authShowJWTCmd = &cobra.Command{
	Use:   "show-jwt",
	Short: "Generate and display the JWT that would be sent",
	Long: `Generate a JWT and display its decoded contents.

This helps verify that:
- The header has the correct algorithm (RS384)
- The payload claims (iss, sub, aud) are correct
- The JWT can be signed with your private key

You can copy the full JWT for manual testing with curl.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if config.CurrentProfile == nil {
			return fmt.Errorf("no profile configured. Run 'fhir-cli config init' first")
		}

		return auth.ShowJWT(config.CurrentProfile)
	},
}

var authStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check authentication status",
	RunE: func(cmd *cobra.Command, args []string) error {
		if config.CurrentProfile == nil {
			return fmt.Errorf("no profile configured. Run 'fhir-cli config init' first")
		}

		profileName := config.CurrentProfile.Name
		if profileName == "" {
			profileName = "default"
		}

		cache, err := auth.LoadTokenFromKeyring(profileName)
		if err != nil {
			fmt.Println("Status: Not authenticated")
			fmt.Println("Run 'fhir-cli auth token' to authenticate")
			return nil
		}

		fmt.Println("Status: Authenticated")
		fmt.Printf("Profile: %s\n", profileName)
		fmt.Printf("Expires: %s\n", cache.ExpiresAt.Format("2006-01-02 15:04:05"))
		fmt.Printf("Scopes: %s\n", cache.Scope)

		return nil
	},
}

var authLogoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Clear cached authentication tokens",
	RunE: func(cmd *cobra.Command, args []string) error {
		profileName := GetProfile()
		if profileName == "default" && config.CurrentConfig != nil {
			profileName = config.CurrentConfig.Default
		}

		if err := auth.ClearTokenCache(profileName); err != nil {
			return fmt.Errorf("failed to clear token cache: %w", err)
		}

		fmt.Printf("Logged out from profile [%s]\n", profileName)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(authCmd)
	authCmd.AddCommand(authTokenCmd)
	authCmd.AddCommand(authStatusCmd)
	authCmd.AddCommand(authLogoutCmd)
	authCmd.AddCommand(authDebugCmd)
	authCmd.AddCommand(authTestKeyCmd)
	authCmd.AddCommand(authShowJWTCmd)

	authTokenCmd.Flags().BoolVar(&forceRefresh, "force", false, "Force token refresh even if cached token is valid")
	authTokenCmd.Flags().BoolVar(&debugMode, "debug", false, "Show debug information instead of getting token")
	authTokenCmd.Flags().BoolVarP(&verboseMode, "verbose", "v", false, "Show detailed request and response")
}
