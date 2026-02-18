package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Profiles map[string]Profile `yaml:"profiles"`
	Default  string             `yaml:"default"`
}

type Profile struct {
	Name         string `yaml:"name"`
	ClientID     string `yaml:"client_id"`
	PrivateKey   string `yaml:"private_key"`
	TokenURL     string `yaml:"token_url"`
	FHIRBaseURL  string `yaml:"fhir_base_url"`
	FHIRVersion  string `yaml:"fhir_version"`
	Scopes       string `yaml:"scopes"`
	OutputFormat string `yaml:"output_format"`
}

type TokenCache struct {
	AccessToken string    `yaml:"access_token" json:"access_token"`
	ExpiresAt   time.Time `yaml:"expires_at" json:"expires_at"`
	Scope       string    `yaml:"scope" json:"scope"`
}

var (
	CurrentConfig  *Config
	CurrentProfile *Profile
	configDir      string
)

func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		configDir = ".fhir-cli"
	} else {
		configDir = filepath.Join(home, ".fhir-cli")
	}
}

func GetConfigDir() string {
	return configDir
}

func GetConfigPath() string {
	return filepath.Join(configDir, "config.yaml")
}

func GetTokenCachePath(profile string) string {
	return filepath.Join(configDir, fmt.Sprintf("token_%s.yaml", profile))
}

func EnsureConfigDir() error {
	return os.MkdirAll(configDir, 0700)
}

// LoadFile reads and parses the config file, populating CurrentConfig.
// It does not resolve a profile. Returns nil if the file does not exist.
func LoadFile() error {
	configPath := GetConfigPath()

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("failed to read config: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}

	CurrentConfig = &cfg
	return nil
}

func Load(profileName string) error {
	if CurrentConfig == nil {
		if err := LoadFile(); err != nil {
			return err
		}
	}

	if CurrentConfig == nil {
		return fmt.Errorf("config file not found. Run 'fhir-cli config init' to create one")
	}

	if profileName == "" || profileName == "default" {
		profileName = CurrentConfig.Default
	}

	profile, exists := CurrentConfig.Profiles[profileName]
	if !exists {
		return fmt.Errorf("profile '%s' not found in config", profileName)
	}

	CurrentProfile = &profile
	return nil
}

func Save(cfg *Config) error {
	if err := EnsureConfigDir(); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	configPath := GetConfigPath()
	if err := os.WriteFile(configPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

func GetDefaultProfile() Profile {
	return Profile{
		Name:         "sandbox",
		FHIRVersion:  "R4",
		FHIRBaseURL:  "https://fhir.epic.com/interconnect-fhir-oauth/api/FHIR/R4",
		TokenURL:     "https://fhir.epic.com/interconnect-fhir-oauth/oauth2/token",
		Scopes:       "system/Patient.read system/Observation.read system/Condition.read",
		OutputFormat: "json",
	}
}
