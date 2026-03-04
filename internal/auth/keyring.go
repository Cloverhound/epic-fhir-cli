package auth

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Cloverhound/epic-fhir-cli/internal/config"
	"github.com/zalando/go-keyring"
	"gopkg.in/yaml.v3"
)

const keyringService = "fhir-cli"

// SaveTokenToKeyring stores a token cache in the OS keyring.
func SaveTokenToKeyring(profile string, cache *config.TokenCache) error {
	data, err := json.Marshal(cache)
	if err != nil {
		return fmt.Errorf("failed to marshal token cache: %w", err)
	}
	return keyring.Set(keyringService, profile, string(data))
}

// LoadTokenFromKeyring loads a token cache from the OS keyring.
// If no keyring entry exists, it attempts to migrate from a legacy YAML file.
func LoadTokenFromKeyring(profile string) (*config.TokenCache, error) {
	secret, err := keyring.Get(keyringService, profile)
	if err == nil {
		var cache config.TokenCache
		if err := json.Unmarshal([]byte(secret), &cache); err != nil {
			return nil, fmt.Errorf("failed to unmarshal token cache: %w", err)
		}
		return &cache, nil
	}

	// Attempt legacy YAML file migration
	cache, legacyErr := loadLegacyTokenCache(profile)
	if legacyErr != nil {
		return nil, err // return original keyring error
	}

	// Migrate to keyring and remove legacy file
	if saveErr := SaveTokenToKeyring(profile, cache); saveErr != nil {
		return cache, nil // return token even if keyring save fails
	}
	legacyPath := config.GetTokenCachePath(profile)
	os.Remove(legacyPath) // best-effort cleanup

	return cache, nil
}

// DeleteTokenFromKeyring removes a token from the OS keyring.
func DeleteTokenFromKeyring(profile string) error {
	err := keyring.Delete(keyringService, profile)
	if err != nil && err != keyring.ErrNotFound {
		return fmt.Errorf("failed to delete token from keyring: %w", err)
	}
	return nil
}

// loadLegacyTokenCache reads a token from the old YAML file format.
func loadLegacyTokenCache(profile string) (*config.TokenCache, error) {
	path := config.GetTokenCachePath(profile)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cache config.TokenCache
	if err := yaml.Unmarshal(data, &cache); err != nil {
		return nil, fmt.Errorf("failed to parse legacy token cache: %w", err)
	}
	return &cache, nil
}
