package cmd

import (
	"fmt"

	"github.com/Cloverhound/epic-fhir-cli/internal/config"
	"github.com/Cloverhound/epic-fhir-cli/internal/fhir"
	"github.com/Cloverhound/epic-fhir-cli/internal/output"
	"github.com/spf13/cobra"
)

var medicationCmd = &cobra.Command{
	Use:   "medication",
	Short: "Medication resource operations",
	Long:  `Manage MedicationRequest resources - prescriptions and medication orders.`,
}

var medicationGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get a medication request by ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]

		client := fhir.NewClient(config.CurrentProfile)
		result, err := client.Get(fhir.ResourceMedicationRequest, id)
		if err != nil {
			return err
		}

		return output.Print(result, output.ParseFormat(GetOutput()))
	},
}

var medicationSearchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for medication requests",
	Long: `Search for medication requests/prescriptions.

Status values:
  - active: The prescription is active
  - on-hold: The prescription is on hold
  - cancelled: The prescription was cancelled
  - completed: The prescription has been completed
  - stopped: The prescription was stopped
  - draft: The prescription is a draft

Intent values:
  - proposal, plan, order, original-order, reflex-order, filler-order, instance-order, option

Examples:
  fhir-cli medication search --patient 123
  fhir-cli medication search --patient 123 --status active
  fhir-cli medication search --patient 123 --intent order`,
	RunE: func(cmd *cobra.Command, args []string) error {
		params := make(map[string]string)

		if patient, _ := cmd.Flags().GetString("patient"); patient != "" {
			params["patient"] = patient
		}
		if status, _ := cmd.Flags().GetString("status"); status != "" {
			params["status"] = status
		}
		if intent, _ := cmd.Flags().GetString("intent"); intent != "" {
			params["intent"] = intent
		}
		if authored, _ := cmd.Flags().GetString("authoredon"); authored != "" {
			params["authoredon"] = authored
		}

		if count, _ := cmd.Flags().GetInt("count"); count > 0 {
			params["_count"] = fmt.Sprintf("%d", count)
		}

		client := fhir.NewClient(config.CurrentProfile)
		bundle, err := client.Search(fhir.ResourceMedicationRequest, params)
		if err != nil {
			return err
		}

		return output.Print(bundle, output.ParseFormat(GetOutput()))
	},
}

var medicationActiveCmd = &cobra.Command{
	Use:   "active <patient-id>",
	Short: "Get active medications for a patient",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		patientID := args[0]
		params := map[string]string{
			"patient": patientID,
			"status":  "active",
		}

		if count, _ := cmd.Flags().GetInt("count"); count > 0 {
			params["_count"] = fmt.Sprintf("%d", count)
		}

		client := fhir.NewClient(config.CurrentProfile)
		bundle, err := client.Search(fhir.ResourceMedicationRequest, params)
		if err != nil {
			return err
		}

		return output.Print(bundle, output.ParseFormat(GetOutput()))
	},
}

// MedicationStatement commands
var medicationStatementCmd = &cobra.Command{
	Use:   "statement",
	Short: "Medication statement operations",
}

var medicationStatementSearchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for medication statements",
	RunE: func(cmd *cobra.Command, args []string) error {
		params := make(map[string]string)

		if patient, _ := cmd.Flags().GetString("patient"); patient != "" {
			params["patient"] = patient
		}
		if status, _ := cmd.Flags().GetString("status"); status != "" {
			params["status"] = status
		}

		if count, _ := cmd.Flags().GetInt("count"); count > 0 {
			params["_count"] = fmt.Sprintf("%d", count)
		}

		client := fhir.NewClient(config.CurrentProfile)
		bundle, err := client.Search(fhir.ResourceMedicationStatement, params)
		if err != nil {
			return err
		}

		return output.Print(bundle, output.ParseFormat(GetOutput()))
	},
}

// MedicationAdministration commands
var medicationAdminCmd = &cobra.Command{
	Use:   "administration",
	Short: "Medication administration operations",
}

var medicationAdminSearchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for medication administrations",
	RunE: func(cmd *cobra.Command, args []string) error {
		params := make(map[string]string)

		if patient, _ := cmd.Flags().GetString("patient"); patient != "" {
			params["patient"] = patient
		}
		if status, _ := cmd.Flags().GetString("status"); status != "" {
			params["status"] = status
		}
		if effectiveTime, _ := cmd.Flags().GetString("effective-time"); effectiveTime != "" {
			params["effective-time"] = effectiveTime
		}

		if count, _ := cmd.Flags().GetInt("count"); count > 0 {
			params["_count"] = fmt.Sprintf("%d", count)
		}

		client := fhir.NewClient(config.CurrentProfile)
		bundle, err := client.Search(fhir.ResourceMedicationAdministration, params)
		if err != nil {
			return err
		}

		return output.Print(bundle, output.ParseFormat(GetOutput()))
	},
}

// MedicationDispense commands
var medicationDispenseCmd = &cobra.Command{
	Use:   "dispense",
	Short: "Medication dispense operations",
}

var medicationDispenseSearchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for medication dispenses",
	RunE: func(cmd *cobra.Command, args []string) error {
		params := make(map[string]string)

		if patient, _ := cmd.Flags().GetString("patient"); patient != "" {
			params["patient"] = patient
		}
		if status, _ := cmd.Flags().GetString("status"); status != "" {
			params["status"] = status
		}

		if count, _ := cmd.Flags().GetInt("count"); count > 0 {
			params["_count"] = fmt.Sprintf("%d", count)
		}

		client := fhir.NewClient(config.CurrentProfile)
		bundle, err := client.Search(fhir.ResourceMedicationDispense, params)
		if err != nil {
			return err
		}

		return output.Print(bundle, output.ParseFormat(GetOutput()))
	},
}

func init() {
	rootCmd.AddCommand(medicationCmd)
	medicationCmd.AddCommand(medicationGetCmd)
	medicationCmd.AddCommand(medicationSearchCmd)
	medicationCmd.AddCommand(medicationActiveCmd)
	medicationCmd.AddCommand(medicationStatementCmd)
	medicationCmd.AddCommand(medicationAdminCmd)
	medicationCmd.AddCommand(medicationDispenseCmd)

	medicationStatementCmd.AddCommand(medicationStatementSearchCmd)
	medicationAdminCmd.AddCommand(medicationAdminSearchCmd)
	medicationDispenseCmd.AddCommand(medicationDispenseSearchCmd)

	// MedicationRequest search flags
	medicationSearchCmd.Flags().String("patient", "", "Patient FHIR ID")
	medicationSearchCmd.Flags().String("status", "", "Status (active, on-hold, cancelled, completed, stopped)")
	medicationSearchCmd.Flags().String("intent", "", "Intent (order, proposal, plan)")
	medicationSearchCmd.Flags().String("authoredon", "", "Authored date range")
	medicationSearchCmd.Flags().Int("count", 20, "Maximum number of results")

	// Active flags
	medicationActiveCmd.Flags().Int("count", 50, "Maximum number of results")

	// Statement search flags
	medicationStatementSearchCmd.Flags().String("patient", "", "Patient FHIR ID")
	medicationStatementSearchCmd.Flags().String("status", "", "Status")
	medicationStatementSearchCmd.Flags().Int("count", 20, "Maximum number of results")

	// Administration search flags
	medicationAdminSearchCmd.Flags().String("patient", "", "Patient FHIR ID")
	medicationAdminSearchCmd.Flags().String("status", "", "Status")
	medicationAdminSearchCmd.Flags().String("effective-time", "", "Effective time range")
	medicationAdminSearchCmd.Flags().Int("count", 20, "Maximum number of results")

	// Dispense search flags
	medicationDispenseSearchCmd.Flags().String("patient", "", "Patient FHIR ID")
	medicationDispenseSearchCmd.Flags().String("status", "", "Status")
	medicationDispenseSearchCmd.Flags().Int("count", 20, "Maximum number of results")
}
