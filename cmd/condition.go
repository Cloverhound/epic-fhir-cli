package cmd

import (
	"fmt"

	"github.com/jbogarin/fhir-cli/internal/config"
	"github.com/jbogarin/fhir-cli/internal/fhir"
	"github.com/jbogarin/fhir-cli/internal/output"
	"github.com/spf13/cobra"
)

var conditionCmd = &cobra.Command{
	Use:   "condition",
	Short: "Condition resource operations",
	Long:  `Manage Condition resources - problems, diagnoses, and health concerns.`,
}

var conditionGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get a condition by ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]

		client := fhir.NewClient(config.CurrentProfile)
		result, err := client.Get(fhir.ResourceCondition, id)
		if err != nil {
			return err
		}

		return output.Print(result, output.ParseFormat(GetOutput()))
	},
}

var conditionSearchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for conditions",
	Long: `Search for conditions/diagnoses using various criteria.

Clinical Status values:
  - active: The condition is active
  - recurrence: The condition has recurred
  - relapse: The condition has relapsed
  - inactive: The condition is inactive
  - remission: The condition is in remission
  - resolved: The condition is resolved

Categories:
  - problem-list-item: Problems maintained on the problem list
  - encounter-diagnosis: Point in time diagnosis

Examples:
  fhir-cli condition search --patient 123
  fhir-cli condition search --patient 123 --clinical-status active
  fhir-cli condition search --patient 123 --category problem-list-item`,
	RunE: func(cmd *cobra.Command, args []string) error {
		params := make(map[string]string)

		if patient, _ := cmd.Flags().GetString("patient"); patient != "" {
			params["patient"] = patient
		}
		if clinicalStatus, _ := cmd.Flags().GetString("clinical-status"); clinicalStatus != "" {
			params["clinical-status"] = clinicalStatus
		}
		if category, _ := cmd.Flags().GetString("category"); category != "" {
			params["category"] = category
		}
		if code, _ := cmd.Flags().GetString("code"); code != "" {
			params["code"] = code
		}
		if onset, _ := cmd.Flags().GetString("onset-date"); onset != "" {
			params["onset-date"] = onset
		}

		if count, _ := cmd.Flags().GetInt("count"); count > 0 {
			params["_count"] = fmt.Sprintf("%d", count)
		}

		client := fhir.NewClient(config.CurrentProfile)
		bundle, err := client.Search(fhir.ResourceCondition, params)
		if err != nil {
			return err
		}

		return output.Print(bundle, output.ParseFormat(GetOutput()))
	},
}

var conditionActiveCmd = &cobra.Command{
	Use:   "active <patient-id>",
	Short: "Get active conditions for a patient",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		patientID := args[0]
		params := map[string]string{
			"patient":         patientID,
			"clinical-status": "active",
		}

		if count, _ := cmd.Flags().GetInt("count"); count > 0 {
			params["_count"] = fmt.Sprintf("%d", count)
		}

		client := fhir.NewClient(config.CurrentProfile)
		bundle, err := client.Search(fhir.ResourceCondition, params)
		if err != nil {
			return err
		}

		return output.Print(bundle, output.ParseFormat(GetOutput()))
	},
}

var conditionProblemsCmd = &cobra.Command{
	Use:   "problems <patient-id>",
	Short: "Get problem list for a patient",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		patientID := args[0]
		params := map[string]string{
			"patient":  patientID,
			"category": "problem-list-item",
		}

		if count, _ := cmd.Flags().GetInt("count"); count > 0 {
			params["_count"] = fmt.Sprintf("%d", count)
		}

		client := fhir.NewClient(config.CurrentProfile)
		bundle, err := client.Search(fhir.ResourceCondition, params)
		if err != nil {
			return err
		}

		return output.Print(bundle, output.ParseFormat(GetOutput()))
	},
}

func init() {
	rootCmd.AddCommand(conditionCmd)
	conditionCmd.AddCommand(conditionGetCmd)
	conditionCmd.AddCommand(conditionSearchCmd)
	conditionCmd.AddCommand(conditionActiveCmd)
	conditionCmd.AddCommand(conditionProblemsCmd)

	// Search flags
	conditionSearchCmd.Flags().String("patient", "", "Patient FHIR ID")
	conditionSearchCmd.Flags().String("clinical-status", "", "Clinical status (active, inactive, resolved)")
	conditionSearchCmd.Flags().String("category", "", "Category (problem-list-item, encounter-diagnosis)")
	conditionSearchCmd.Flags().String("code", "", "Condition code (ICD-10, SNOMED)")
	conditionSearchCmd.Flags().String("onset-date", "", "Onset date range")
	conditionSearchCmd.Flags().Int("count", 20, "Maximum number of results")

	// Active flags
	conditionActiveCmd.Flags().Int("count", 20, "Maximum number of results")

	// Problems flags
	conditionProblemsCmd.Flags().Int("count", 20, "Maximum number of results")
}
