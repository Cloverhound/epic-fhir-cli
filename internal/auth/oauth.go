package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/Cloverhound/epic-fhir-cli/internal/config"
)

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
}

type OAuthError struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

func GetAccessToken(profile *config.Profile, forceRefresh bool) (string, error) {
	profileName := profile.Name
	if profileName == "" {
		profileName = "default"
	}

	// Check token cache first
	if !forceRefresh {
		cache, err := LoadTokenFromKeyring(profileName)
		if err == nil && cache.ExpiresAt.After(time.Now().Add(5*time.Minute)) {
			return cache.AccessToken, nil
		}
	}

	// Generate new JWT assertion
	jwtAssertion, err := GenerateJWT(profile.ClientID, profile.TokenURL, profile.PrivateKey)
	if err != nil {
		return "", fmt.Errorf("failed to generate JWT: %w", err)
	}

	// Exchange JWT for access token
	// Note: For Epic backend OAuth2, scopes are determined by app registration, not sent in request
	tokenResp, err := exchangeJWTForToken(profile.TokenURL, jwtAssertion)
	if err != nil {
		return "", err
	}

	// Cache the token
	cache := &config.TokenCache{
		AccessToken: tokenResp.AccessToken,
		ExpiresAt:   time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second),
		Scope:       tokenResp.Scope,
	}
	if err := SaveTokenToKeyring(profileName, cache); err != nil {
		// Log warning but don't fail
		fmt.Printf("Warning: failed to cache token: %v\n", err)
	}

	return tokenResp.AccessToken, nil
}

// VerboseMode enables detailed request/response logging
var VerboseMode bool

func exchangeJWTForToken(tokenURL, jwtAssertion string) (*TokenResponse, error) {
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_assertion_type", "urn:ietf:params:oauth:client-assertion-type:jwt-bearer")
	data.Set("client_assertion", jwtAssertion)

	if VerboseMode {
		fmt.Println("\n=== Token Request ===")
		fmt.Printf("POST %s\n", tokenURL)
		fmt.Println("Content-Type: application/x-www-form-urlencoded")
		fmt.Printf("Body:\n")
		fmt.Printf("  grant_type=%s\n", data.Get("grant_type"))
		fmt.Printf("  client_assertion_type=%s\n", data.Get("client_assertion_type"))
		fmt.Printf("  client_assertion=%s...\n", jwtAssertion[:min(50, len(jwtAssertion))])
	}

	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create token request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("token request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read token response: %w", err)
	}

	if VerboseMode {
		fmt.Println("\n=== Token Response ===")
		fmt.Printf("Status: %d %s\n", resp.StatusCode, resp.Status)
		fmt.Println("Headers:")
		for k, v := range resp.Header {
			fmt.Printf("  %s: %s\n", k, strings.Join(v, ", "))
		}
		fmt.Println("Body:")
		// Pretty print JSON if possible
		var prettyJSON map[string]interface{}
		if json.Unmarshal(body, &prettyJSON) == nil {
			prettyBytes, _ := json.MarshalIndent(prettyJSON, "  ", "  ")
			fmt.Printf("  %s\n", string(prettyBytes))
		} else {
			fmt.Printf("  %s\n", string(body))
		}
	}

	if resp.StatusCode != http.StatusOK {
		var oauthErr OAuthError
		if json.Unmarshal(body, &oauthErr) == nil && oauthErr.Error != "" {
			return nil, fmt.Errorf("OAuth error: %s - %s", oauthErr.Error, oauthErr.ErrorDescription)
		}
		return nil, fmt.Errorf("token request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var tokenResp TokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, fmt.Errorf("failed to parse token response: %w", err)
	}

	return &tokenResp, nil
}


func ClearTokenCache(profileName string) error {
	if err := DeleteTokenFromKeyring(profileName); err != nil {
		return err
	}
	// Also remove legacy YAML file if it exists
	legacyPath := config.GetTokenCachePath(profileName)
	if err := os.Remove(legacyPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove legacy token cache: %w", err)
	}
	return nil
}
