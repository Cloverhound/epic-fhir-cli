package cmd

import (
	"fmt"

	"github.com/Cloverhound/epic-fhir-cli/internal/config"
	"github.com/Cloverhound/epic-fhir-cli/internal/fhir"
	"github.com/Cloverhound/epic-fhir-cli/internal/output"
	"github.com/spf13/cobra"
)

var encounterCmd = &cobra.Command{
	Use:   "encounter",
	Short: "Encounter resource operations",
	Long:  `Manage Encounter resources - healthcare visits and admissions.`,
}

var encounterGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get an encounter by ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]

		client := fhir.NewClient(config.CurrentProfile)
		result, err := client.Get(fhir.ResourceEncounter, id)
		if err != nil {
			return err
		}

		return output.Print(result, output.ParseFormat(GetOutput()))
	},
}

var encounterSearchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for encounters",
	Long: `Search for encounter records.

Status values:
  - planned, arrived, triaged, in-progress, onleave, finished, cancelled

Class values:
  - AMB (ambulatory), EMER (emergency), IMP (inpatient), OBSENC (observation)

Examples:
  fhir-cli encounter search --patient 123
  fhir-cli encounter search --patient 123 --status finished
  fhir-cli encounter search --patient 123 --class IMP
  fhir-cli encounter search --patient 123 --date ge2024-01-01`,
	RunE: func(cmd *cobra.Command, args []string) error {
		params := make(map[string]string)

		if patient, _ := cmd.Flags().GetString("patient"); patient != "" {
			params["patient"] = patient
		}
		if status, _ := cmd.Flags().GetString("status"); status != "" {
			params["status"] = status
		}
		if class, _ := cmd.Flags().GetString("class"); class != "" {
			params["class"] = class
		}
		if date, _ := cmd.Flags().GetString("date"); date != "" {
			params["date"] = date
		}
		if encType, _ := cmd.Flags().GetString("type"); encType != "" {
			params["type"] = encType
		}

		if count, _ := cmd.Flags().GetInt("count"); count > 0 {
			params["_count"] = fmt.Sprintf("%d", count)
		}

		client := fhir.NewClient(config.CurrentProfile)
		bundle, err := client.Search(fhir.ResourceEncounter, params)
		if err != nil {
			return err
		}

		return output.Print(bundle, output.ParseFormat(GetOutput()))
	},
}

var encounterListCmd = &cobra.Command{
	Use:   "list <patient-id>",
	Short: "List encounters for a patient",
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
		bundle, err := client.Search(fhir.ResourceEncounter, params)
		if err != nil {
			return err
		}

		return output.Print(bundle, output.ParseFormat(GetOutput()))
	},
}

func init() {
	rootCmd.AddCommand(encounterCmd)
	encounterCmd.AddCommand(encounterGetCmd)
	encounterCmd.AddCommand(encounterSearchCmd)
	encounterCmd.AddCommand(encounterListCmd)

	// Search flags
	encounterSearchCmd.Flags().String("patient", "", "Patient FHIR ID")
	encounterSearchCmd.Flags().String("status", "", "Status (planned, in-progress, finished, cancelled)")
	encounterSearchCmd.Flags().String("class", "", "Class (AMB, EMER, IMP, OBSENC)")
	encounterSearchCmd.Flags().String("date", "", "Date range")
	encounterSearchCmd.Flags().String("type", "", "Encounter type")
	encounterSearchCmd.Flags().Int("count", 20, "Maximum number of results")

	// List flags
	encounterListCmd.Flags().Int("count", 20, "Maximum number of results")
}
