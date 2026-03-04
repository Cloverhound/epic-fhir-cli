package auth

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Cloverhound/epic-fhir-cli/internal/config"
)

// ShowJWT generates and displays the full decoded JWT that would be sent
func ShowJWT(profile *config.Profile) error {
	fmt.Println("=== JWT Debug ===")
	fmt.Println()

	// Generate the JWT
	jwtToken, err := GenerateJWT(profile.ClientID, profile.TokenURL, profile.PrivateKey)
	if err != nil {
		return fmt.Errorf("failed to generate JWT: %w", err)
	}

	parts := strings.Split(jwtToken, ".")
	if len(parts) != 3 {
		return fmt.Errorf("invalid JWT format")
	}

	// Decode header
	fmt.Println("Header:")
	if headerJSON, err := base64.RawURLEncoding.DecodeString(parts[0]); err == nil {
		var header map[string]interface{}
		json.Unmarshal(headerJSON, &header)
		prettyHeader, _ := json.MarshalIndent(header, "  ", "  ")
		fmt.Printf("  %s\n", string(prettyHeader))
	}

	// Decode payload
	fmt.Println("\nPayload (Claims):")
	if payloadJSON, err := base64.RawURLEncoding.DecodeString(parts[1]); err == nil {
		var payload map[string]interface{}
		json.Unmarshal(payloadJSON, &payload)
		prettyPayload, _ := json.MarshalIndent(payload, "  ", "  ")
		fmt.Printf("  %s\n", string(prettyPayload))
	}

	fmt.Println("\nFull JWT (for manual testing):")
	fmt.Println(jwtToken)

	fmt.Println("\n=== Verification Checklist ===")
	fmt.Println("1. 'iss' (issuer) must match your registered client_id exactly")
	fmt.Println("2. 'sub' (subject) must match your registered client_id exactly")
	fmt.Println("3. 'aud' (audience) must match the token endpoint URL exactly")
	fmt.Println("4. 'alg' in header must be RS384")
	fmt.Println("5. The private key must match the public key registered with Epic")

	return nil
}

// DebugAuth prints debug information about the authentication configuration
func DebugAuth(profile *config.Profile) error {
	fmt.Println("=== Authentication Debug Info ===")
	fmt.Println()

	// Profile info
	fmt.Println("Profile Configuration:")
	fmt.Printf("  Name:          %s\n", profile.Name)
	fmt.Printf("  Client ID:     %s\n", profile.ClientID)
	fmt.Printf("  Token URL:     %s\n", profile.TokenURL)
	fmt.Printf("  FHIR Base URL: %s\n", profile.FHIRBaseURL)
	fmt.Printf("  Scopes:        %s\n", profile.Scopes)
	fmt.Println()

	// Private key info
	fmt.Println("Private Key:")
	fmt.Printf("  Source: %s\n", profile.PrivateKey)
	if strings.HasPrefix(profile.PrivateKey, "env:") {
		envVar := strings.TrimPrefix(profile.PrivateKey, "env:")
		envValue := os.Getenv(envVar)
		if envValue == "" {
			fmt.Printf("  Status: ❌ Environment variable %s is NOT SET\n", envVar)
		} else {
			fmt.Printf("  Status: ✓ Environment variable %s is set (%d bytes)\n", envVar, len(envValue))
			// Show first line of the key to verify format
			lines := strings.Split(envValue, "\n")
			if len(lines) > 0 {
				fmt.Printf("  Format: %s\n", lines[0])
			}
		}
	} else {
		if _, err := os.Stat(profile.PrivateKey); os.IsNotExist(err) {
			fmt.Printf("  Status: ❌ File does not exist\n")
		} else {
			fmt.Printf("  Status: ✓ File exists\n")
		}
	}

	// Try to load the key
	fmt.Println()
	fmt.Println("Key Validation:")
	key, err := loadPrivateKey(profile.PrivateKey)
	if err != nil {
		fmt.Printf("  ❌ Failed to load key: %v\n", err)
	} else {
		fmt.Printf("  ✓ Key loaded successfully\n")
		fmt.Printf("  Key size: %d bits\n", key.N.BitLen())
	}

	// Show JWT claims that would be generated
	fmt.Println()
	fmt.Println("JWT Claims (what will be sent):")
	now := time.Now()
	fmt.Printf("  iss (issuer):   %s\n", profile.ClientID)
	fmt.Printf("  sub (subject):  %s\n", profile.ClientID)
	fmt.Printf("  aud (audience): %s\n", profile.TokenURL)
	fmt.Printf("  iat (issued):   %s\n", now.Format(time.RFC3339))
	fmt.Printf("  exp (expires):  %s\n", now.Add(5*time.Minute).Format(time.RFC3339))
	fmt.Printf("  jti (jwt id):   <random UUID>\n")
	fmt.Printf("  Algorithm:      RS384\n")

	// If key loaded, try to generate a JWT
	if key != nil {
		fmt.Println()
		fmt.Println("JWT Generation Test:")
		jwtToken, err := GenerateJWT(profile.ClientID, profile.TokenURL, profile.PrivateKey)
		if err != nil {
			fmt.Printf("  ❌ Failed to generate JWT: %v\n", err)
		} else {
			fmt.Printf("  ✓ JWT generated successfully\n")
			// Show JWT parts
			parts := strings.Split(jwtToken, ".")
			if len(parts) == 3 {
				fmt.Printf("  Header:  %s...\n", parts[0][:min(50, len(parts[0]))])
				fmt.Printf("  Payload: %s...\n", parts[1][:min(50, len(parts[1]))])
				fmt.Printf("  Signature: %s...\n", parts[2][:min(30, len(parts[2]))])

				// Decode and show payload
				if payload, err := base64.RawURLEncoding.DecodeString(parts[1]); err == nil {
					var claims map[string]interface{}
					if json.Unmarshal(payload, &claims) == nil {
						fmt.Println()
						fmt.Println("  Decoded payload:")
						for k, v := range claims {
							fmt.Printf("    %s: %v\n", k, v)
						}
					}
				}
			}
		}
	}

	// OAuth request info
	fmt.Println()
	fmt.Println("OAuth Token Request (what will be sent):")
	fmt.Printf("  POST %s\n", profile.TokenURL)
	fmt.Printf("  Content-Type: application/x-www-form-urlencoded\n")
	fmt.Printf("  Body:\n")
	fmt.Printf("    grant_type=client_credentials\n")
	fmt.Printf("    client_assertion_type=urn:ietf:params:oauth:client-assertion-type:jwt-bearer\n")
	fmt.Printf("    client_assertion=<signed JWT>\n")
	fmt.Println()
	fmt.Println("  Note: Scopes are NOT sent in the request.")
	fmt.Println("  For Epic backend OAuth2, scopes are determined by your app registration.")

	// Common issues
	fmt.Println()
	fmt.Println("=== Common Issues ===")
	fmt.Println()
	fmt.Println("unauthorized_client errors can be caused by:")
	fmt.Println("  1. Wrong client_id - must match your Epic app registration exactly")
	fmt.Println("  2. Wrong token_url - verify it's correct for your Epic environment")
	fmt.Println("  3. Wrong private key - must match the public key registered with Epic")
	fmt.Println("  4. Key format issues - ensure it's a valid RSA private key in PEM format")
	fmt.Println("  5. App not approved - your app may need Epic approval")
	fmt.Println("  6. Wrong audience - the 'aud' claim must match the token URL")
	fmt.Println()
	fmt.Println("For Epic sandbox, verify your app is registered at: https://fhir.epic.com/")

	return nil
}

// TestPrivateKey attempts to load and validate a private key
func TestPrivateKey(source string) error {
	key, err := loadPrivateKey(source)
	if err != nil {
		return err
	}

	// Validate key can be used for signing
	if key.N.BitLen() < 2048 {
		return fmt.Errorf("key size %d bits is too small (minimum 2048)", key.N.BitLen())
	}

	return nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
