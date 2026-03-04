package cmd

import (
	"fmt"

	"github.com/Cloverhound/epic-fhir-cli/internal/config"
	"github.com/Cloverhound/epic-fhir-cli/internal/fhir"
	"github.com/Cloverhound/epic-fhir-cli/internal/output"
	"github.com/spf13/cobra"
)

var observationCmd = &cobra.Command{
	Use:   "observation",
	Short: "Observation resource operations",
	Long:  `Manage Observation resources - measurements and simple assertions.`,
}

var observationGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get an observation by ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]

		client := fhir.NewClient(config.CurrentProfile)
		result, err := client.Get(fhir.ResourceObservation, id)
		if err != nil {
			return err
		}

		return output.Print(result, output.ParseFormat(GetOutput()))
	},
}

var observationSearchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for observations",
	Long: `Search for observations using various criteria.

Categories:
  - vital-signs: Blood pressure, heart rate, temperature, etc.
  - laboratory: Lab results
  - social-history: Smoking status, alcohol use, etc.
  - imaging: Imaging observations
  - exam: Physical exam findings

Examples:
  fhir-cli observation search --patient 123
  fhir-cli observation search --patient 123 --category vital-signs
  fhir-cli observation search --patient 123 --code 8867-4 # Heart rate
  fhir-cli observation search --patient 123 --date ge2024-01-01`,
	RunE: func(cmd *cobra.Command, args []string) error {
		params := make(map[string]string)

		if patient, _ := cmd.Flags().GetString("patient"); patient != "" {
			params["patient"] = patient
		}
		if category, _ := cmd.Flags().GetString("category"); category != "" {
			params["category"] = category
		}
		if code, _ := cmd.Flags().GetString("code"); code != "" {
			params["code"] = code
		}
		if date, _ := cmd.Flags().GetString("date"); date != "" {
			params["date"] = date
		}
		if status, _ := cmd.Flags().GetString("status"); status != "" {
			params["status"] = status
		}

		if count, _ := cmd.Flags().GetInt("count"); count > 0 {
			params["_count"] = fmt.Sprintf("%d", count)
		}

		client := fhir.NewClient(config.CurrentProfile)
		bundle, err := client.Search(fhir.ResourceObservation, params)
		if err != nil {
			return err
		}

		return output.Print(bundle, output.ParseFormat(GetOutput()))
	},
}

var observationVitalsCmd = &cobra.Command{
	Use:   "vitals <patient-id>",
	Short: "Get vital signs for a patient",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		patientID := args[0]
		params := map[string]string{
			"patient":  patientID,
			"category": "vital-signs",
		}

		if count, _ := cmd.Flags().GetInt("count"); count > 0 {
			params["_count"] = fmt.Sprintf("%d", count)
		}

		if date, _ := cmd.Flags().GetString("date"); date != "" {
			params["date"] = date
		}

		client := fhir.NewClient(config.CurrentProfile)
		bundle, err := client.Search(fhir.ResourceObservation, params)
		if err != nil {
			return err
		}

		return output.Print(bundle, output.ParseFormat(GetOutput()))
	},
}

var observationLabsCmd = &cobra.Command{
	Use:   "labs <patient-id>",
	Short: "Get laboratory results for a patient",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		patientID := args[0]
		params := map[string]string{
			"patient":  patientID,
			"category": "laboratory",
		}

		if count, _ := cmd.Flags().GetInt("count"); count > 0 {
			params["_count"] = fmt.Sprintf("%d", count)
		}

		if date, _ := cmd.Flags().GetString("date"); date != "" {
			params["date"] = date
		}

		client := fhir.NewClient(config.CurrentProfile)
		bundle, err := client.Search(fhir.ResourceObservation, params)
		if err != nil {
			return err
		}

		return output.Print(bundle, output.ParseFormat(GetOutput()))
	},
}

func init() {
	rootCmd.AddCommand(observationCmd)
	observationCmd.AddCommand(observationGetCmd)
	observationCmd.AddCommand(observationSearchCmd)
	observationCmd.AddCommand(observationVitalsCmd)
	observationCmd.AddCommand(observationLabsCmd)

	// Search flags
	observationSearchCmd.Flags().String("patient", "", "Patient FHIR ID (required)")
	observationSearchCmd.Flags().String("category", "", "Category (vital-signs, laboratory, social-history, etc.)")
	observationSearchCmd.Flags().String("code", "", "LOINC or other code")
	observationSearchCmd.Flags().String("date", "", "Date range (e.g., ge2024-01-01, le2024-12-31)")
	observationSearchCmd.Flags().String("status", "", "Status (registered, preliminary, final, amended)")
	observationSearchCmd.Flags().Int("count", 20, "Maximum number of results")

	// Vitals flags
	observationVitalsCmd.Flags().Int("count", 20, "Maximum number of results")
	observationVitalsCmd.Flags().String("date", "", "Date range")

	// Labs flags
	observationLabsCmd.Flags().Int("count", 20, "Maximum number of results")
	observationLabsCmd.Flags().String("date", "", "Date range")
}
