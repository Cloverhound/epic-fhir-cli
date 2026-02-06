package cmd

import (
	"fmt"

	"github.com/jbogarin/fhir-cli/internal/config"
	"github.com/jbogarin/fhir-cli/internal/fhir"
	"github.com/jbogarin/fhir-cli/internal/output"
	"github.com/spf13/cobra"
)

var procedureCmd = &cobra.Command{
	Use:   "procedure",
	Short: "Procedure resource operations",
	Long:  `Manage Procedure resources - actions performed on or for a patient.`,
}

var procedureGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get a procedure by ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]

		client := fhir.NewClient(config.CurrentProfile)
		result, err := client.Get(fhir.ResourceProcedure, id)
		if err != nil {
			return err
		}

		return output.Print(result, output.ParseFormat(GetOutput()))
	},
}

var procedureSearchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for procedures",
	Long: `Search for procedure records.

Status values:
  - preparation, in-progress, not-done, on-hold, stopped, completed, entered-in-error, unknown

Examples:
  fhir-cli procedure search --patient 123
  fhir-cli procedure search --patient 123 --status completed
  fhir-cli procedure search --patient 123 --date ge2024-01-01`,
	RunE: func(cmd *cobra.Command, args []string) error {
		params := make(map[string]string)

		if patient, _ := cmd.Flags().GetString("patient"); patient != "" {
			params["patient"] = patient
		}
		if status, _ := cmd.Flags().GetString("status"); status != "" {
			params["status"] = status
		}
		if code, _ := cmd.Flags().GetString("code"); code != "" {
			params["code"] = code
		}
		if date, _ := cmd.Flags().GetString("date"); date != "" {
			params["date"] = date
		}
		if category, _ := cmd.Flags().GetString("category"); category != "" {
			params["category"] = category
		}

		if count, _ := cmd.Flags().GetInt("count"); count > 0 {
			params["_count"] = fmt.Sprintf("%d", count)
		}

		client := fhir.NewClient(config.CurrentProfile)
		bundle, err := client.Search(fhir.ResourceProcedure, params)
		if err != nil {
			return err
		}

		return output.Print(bundle, output.ParseFormat(GetOutput()))
	},
}

var procedureListCmd = &cobra.Command{
	Use:   "list <patient-id>",
	Short: "List procedures for a patient",
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
		bundle, err := client.Search(fhir.ResourceProcedure, params)
		if err != nil {
			return err
		}

		return output.Print(bundle, output.ParseFormat(GetOutput()))
	},
}

func init() {
	rootCmd.AddCommand(procedureCmd)
	procedureCmd.AddCommand(procedureGetCmd)
	procedureCmd.AddCommand(procedureSearchCmd)
	procedureCmd.AddCommand(procedureListCmd)

	// Search flags
	procedureSearchCmd.Flags().String("patient", "", "Patient FHIR ID")
	procedureSearchCmd.Flags().String("status", "", "Status (completed, in-progress, not-done)")
	procedureSearchCmd.Flags().String("code", "", "Procedure code (CPT, SNOMED)")
	procedureSearchCmd.Flags().String("date", "", "Date range")
	procedureSearchCmd.Flags().String("category", "", "Category")
	procedureSearchCmd.Flags().Int("count", 20, "Maximum number of results")

	// List flags
	procedureListCmd.Flags().Int("count", 20, "Maximum number of results")
}
