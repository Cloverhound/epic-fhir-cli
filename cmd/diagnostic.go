package cmd

import (
	"fmt"

	"github.com/Cloverhound/epic-fhir-cli/internal/config"
	"github.com/Cloverhound/epic-fhir-cli/internal/fhir"
	"github.com/Cloverhound/epic-fhir-cli/internal/output"
	"github.com/spf13/cobra"
)

var diagnosticCmd = &cobra.Command{
	Use:     "diagnostic",
	Aliases: []string{"dx"},
	Short:   "DiagnosticReport resource operations",
	Long:    `Manage DiagnosticReport resources - lab results, imaging reports, pathology.`,
}

var diagnosticGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get a diagnostic report by ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]

		client := fhir.NewClient(config.CurrentProfile)
		result, err := client.Get(fhir.ResourceDiagnosticReport, id)
		if err != nil {
			return err
		}

		return output.Print(result, output.ParseFormat(GetOutput()))
	},
}

var diagnosticSearchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for diagnostic reports",
	Long: `Search for diagnostic reports.

Categories:
  - LAB: Laboratory
  - RAD: Radiology/Imaging
  - PATH: Pathology
  - CARD: Cardiology

Status values:
  - registered, partial, preliminary, final, amended, corrected, appended, cancelled, entered-in-error

Examples:
  fhir-cli diagnostic search --patient 123
  fhir-cli diagnostic search --patient 123 --category LAB
  fhir-cli diagnostic search --patient 123 --status final`,
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
		if status, _ := cmd.Flags().GetString("status"); status != "" {
			params["status"] = status
		}
		if date, _ := cmd.Flags().GetString("date"); date != "" {
			params["date"] = date
		}

		if count, _ := cmd.Flags().GetInt("count"); count > 0 {
			params["_count"] = fmt.Sprintf("%d", count)
		}

		client := fhir.NewClient(config.CurrentProfile)
		bundle, err := client.Search(fhir.ResourceDiagnosticReport, params)
		if err != nil {
			return err
		}

		return output.Print(bundle, output.ParseFormat(GetOutput()))
	},
}

var diagnosticLabsCmd = &cobra.Command{
	Use:   "labs <patient-id>",
	Short: "Get lab reports for a patient",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		patientID := args[0]
		params := map[string]string{
			"patient":  patientID,
			"category": "LAB",
		}

		if count, _ := cmd.Flags().GetInt("count"); count > 0 {
			params["_count"] = fmt.Sprintf("%d", count)
		}

		client := fhir.NewClient(config.CurrentProfile)
		bundle, err := client.Search(fhir.ResourceDiagnosticReport, params)
		if err != nil {
			return err
		}

		return output.Print(bundle, output.ParseFormat(GetOutput()))
	},
}

var diagnosticImagingCmd = &cobra.Command{
	Use:   "imaging <patient-id>",
	Short: "Get imaging reports for a patient",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		patientID := args[0]
		params := map[string]string{
			"patient":  patientID,
			"category": "RAD",
		}

		if count, _ := cmd.Flags().GetInt("count"); count > 0 {
			params["_count"] = fmt.Sprintf("%d", count)
		}

		client := fhir.NewClient(config.CurrentProfile)
		bundle, err := client.Search(fhir.ResourceDiagnosticReport, params)
		if err != nil {
			return err
		}

		return output.Print(bundle, output.ParseFormat(GetOutput()))
	},
}

func init() {
	rootCmd.AddCommand(diagnosticCmd)
	diagnosticCmd.AddCommand(diagnosticGetCmd)
	diagnosticCmd.AddCommand(diagnosticSearchCmd)
	diagnosticCmd.AddCommand(diagnosticLabsCmd)
	diagnosticCmd.AddCommand(diagnosticImagingCmd)

	// Search flags
	diagnosticSearchCmd.Flags().String("patient", "", "Patient FHIR ID")
	diagnosticSearchCmd.Flags().String("category", "", "Category (LAB, RAD, PATH, CARD)")
	diagnosticSearchCmd.Flags().String("code", "", "Report code")
	diagnosticSearchCmd.Flags().String("status", "", "Status (final, preliminary, registered)")
	diagnosticSearchCmd.Flags().String("date", "", "Date range")
	diagnosticSearchCmd.Flags().Int("count", 20, "Maximum number of results")

	// Labs flags
	diagnosticLabsCmd.Flags().Int("count", 20, "Maximum number of results")

	// Imaging flags
	diagnosticImagingCmd.Flags().Int("count", 20, "Maximum number of results")
}
