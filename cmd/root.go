package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Cloverhound/epic-fhir-cli/internal/auth"
	"github.com/Cloverhound/epic-fhir-cli/internal/config"
	"github.com/Cloverhound/epic-fhir-cli/internal/fhir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile      string
	outputFormat string
	profile      string
	verbose      bool
	Version      = "0.1.0"
	BuildDate    = "unknown"
)

var rootCmd = &cobra.Command{
	Use:   "fhir-cli",
	Short: "A CLI for querying Epic FHIR APIs",
	Long: `fhir-cli is a command-line interface for interacting with Epic's FHIR R4 APIs.

It supports OAuth2 backend authentication with JWT and provides access to all
FHIR resources including Patient, Observation, Condition, MedicationRequest,
and more.

Configure your credentials with 'fhir-cli config init' to get started.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Set verbose mode for all packages
		auth.VerboseMode = verbose
		fhir.VerboseMode = verbose

		// Skip config loading for commands that don't need it
		switch cmd.Name() {
		case "version", "update":
			return nil
		case "init":
			if cmd.Parent() != nil && cmd.Parent().Name() == "config" {
				return nil
			}
		}
		return config.Load(profile)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.fhir-cli/config.yaml)")
	rootCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "json", "output format: json, table, yaml")
	rootCmd.PersistentFlags().StringVarP(&profile, "profile", "p", "default", "configuration profile to use")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output (show requests/responses)")

	viper.BindPFlag("output", rootCmd.PersistentFlags().Lookup("output"))
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		configDir := filepath.Join(home, ".fhir-cli")
		viper.AddConfigPath(configDir)
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
	}

	viper.SetEnvPrefix("FHIR")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		// Config file not found is okay for init command
	}
}

func GetOutput() string {
	return outputFormat
}

func GetProfile() string {
	return profile
}

func IsVerbose() bool {
	return verbose
}
