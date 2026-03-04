package cmd

import (
	"fmt"

	"github.com/Cloverhound/epic-fhir-cli/internal/config"
	"github.com/Cloverhound/epic-fhir-cli/internal/fhir"
	"github.com/Cloverhound/epic-fhir-cli/internal/output"
	"github.com/spf13/cobra"
)

var immunizationCmd = &cobra.Command{
	Use:     "immunization",
	Aliases: []string{"vaccine", "imm"},
	Short:   "Immunization resource operations",
	Long:    `Manage Immunization resources - vaccination records.`,
}

var immunizationGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get an immunization record by ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]

		client := fhir.NewClient(config.CurrentProfile)
		result, err := client.Get(fhir.ResourceImmunization, id)
		if err != nil {
			return err
		}

		return output.Print(result, output.ParseFormat(GetOutput()))
	},
}

var immunizationSearchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for immunization records",
	Long: `Search for immunization records.

Status values:
  - completed: The immunization was given
  - entered-in-error: The record was entered in error
  - not-done: The immunization was not given

Examples:
  fhir-cli immunization search --patient 123
  fhir-cli immunization search --patient 123 --status completed
  fhir-cli immunization search --patient 123 --date ge2024-01-01`,
	RunE: func(cmd *cobra.Command, args []string) error {
		params := make(map[string]string)

		if patient, _ := cmd.Flags().GetString("patient"); patient != "" {
			params["patient"] = patient
		}
		if status, _ := cmd.Flags().GetString("status"); status != "" {
			params["status"] = status
		}
		if vaccineCode, _ := cmd.Flags().GetString("vaccine-code"); vaccineCode != "" {
			params["vaccine-code"] = vaccineCode
		}
		if date, _ := cmd.Flags().GetString("date"); date != "" {
			params["date"] = date
		}

		if count, _ := cmd.Flags().GetInt("count"); count > 0 {
			params["_count"] = fmt.Sprintf("%d", count)
		}

		client := fhir.NewClient(config.CurrentProfile)
		bundle, err := client.Search(fhir.ResourceImmunization, params)
		if err != nil {
			return err
		}

		return output.Print(bundle, output.ParseFormat(GetOutput()))
	},
}

var immunizationListCmd = &cobra.Command{
	Use:   "list <patient-id>",
	Short: "List immunizations for a patient",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		patientID := args[0]
		params := map[string]string{
			"patient": patientID,
		}

		if count, _ := cmd.Flags().GetInt("count"); count > 0 {
			params["_count"] = fmt.Sprintf("%d", count)
		}

		client := fhir.NewClient(config.CurrentProfile)
		bundle, err := client.Search(fhir.ResourceImmunization, params)
		if err != nil {
			return err
		}

		return output.Print(bundle, output.ParseFormat(GetOutput()))
	},
}

func init() {
	rootCmd.AddCommand(immunizationCmd)
	immunizationCmd.AddCommand(immunizationGetCmd)
	immunizationCmd.AddCommand(immunizationSearchCmd)
	immunizationCmd.AddCommand(immunizationListCmd)

	// Search flags
	immunizationSearchCmd.Flags().String("patient", "", "Patient FHIR ID")
	immunizationSearchCmd.Flags().String("status", "", "Status (completed, not-done, entered-in-error)")
	immunizationSearchCmd.Flags().String("vaccine-code", "", "Vaccine code (CVX)")
	immunizationSearchCmd.Flags().String("date", "", "Date range")
	immunizationSearchCmd.Flags().Int("count", 50, "Maximum number of results")

	// List flags
	immunizationListCmd.Flags().Int("count", 50, "Maximum number of results")
}
