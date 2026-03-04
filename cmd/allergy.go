package cmd

import (
	"fmt"

	"github.com/Cloverhound/epic-fhir-cli/internal/config"
	"github.com/Cloverhound/epic-fhir-cli/internal/fhir"
	"github.com/Cloverhound/epic-fhir-cli/internal/output"
	"github.com/spf13/cobra"
)

var allergyCmd = &cobra.Command{
	Use:     "allergy",
	Aliases: []string{"allergies"},
	Short:   "AllergyIntolerance resource operations",
	Long:    `Manage AllergyIntolerance resources - allergies and adverse reactions.`,
}

var allergyGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get an allergy by ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]

		client := fhir.NewClient(config.CurrentProfile)
		result, err := client.Get(fhir.ResourceAllergyIntolerance, id)
		if err != nil {
			return err
		}

		return output.Print(result, output.ParseFormat(GetOutput()))
	},
}

var allergySearchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for allergies",
	Long: `Search for allergy/intolerance records.

Clinical Status values:
  - active: The allergy is active
  - inactive: The allergy is inactive
  - resolved: The allergy has resolved

Criticality values:
  - low: Low risk
  - high: High risk
  - unable-to-assess: Unable to assess

Type values:
  - allergy: Allergic reaction
  - intolerance: Non-allergic intolerance

Examples:
  fhir-cli allergy search --patient 123
  fhir-cli allergy search --patient 123 --clinical-status active
  fhir-cli allergy search --patient 123 --criticality high`,
	RunE: func(cmd *cobra.Command, args []string) error {
		params := make(map[string]string)

		if patient, _ := cmd.Flags().GetString("patient"); patient != "" {
			params["patient"] = patient
		}
		if clinicalStatus, _ := cmd.Flags().GetString("clinical-status"); clinicalStatus != "" {
			params["clinical-status"] = clinicalStatus
		}
		if criticality, _ := cmd.Flags().GetString("criticality"); criticality != "" {
			params["criticality"] = criticality
		}
		if allergyType, _ := cmd.Flags().GetString("type"); allergyType != "" {
			params["type"] = allergyType
		}
		if category, _ := cmd.Flags().GetString("category"); category != "" {
			params["category"] = category
		}

		if count, _ := cmd.Flags().GetInt("count"); count > 0 {
			params["_count"] = fmt.Sprintf("%d", count)
		}

		client := fhir.NewClient(config.CurrentProfile)
		bundle, err := client.Search(fhir.ResourceAllergyIntolerance, params)
		if err != nil {
			return err
		}

		return output.Print(bundle, output.ParseFormat(GetOutput()))
	},
}

var allergyListCmd = &cobra.Command{
	Use:   "list <patient-id>",
	Short: "List all allergies for a patient",
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
		bundle, err := client.Search(fhir.ResourceAllergyIntolerance, params)
		if err != nil {
			return err
		}

		return output.Print(bundle, output.ParseFormat(GetOutput()))
	},
}

func init() {
	rootCmd.AddCommand(allergyCmd)
	allergyCmd.AddCommand(allergyGetCmd)
	allergyCmd.AddCommand(allergySearchCmd)
	allergyCmd.AddCommand(allergyListCmd)

	// Search flags
	allergySearchCmd.Flags().String("patient", "", "Patient FHIR ID")
	allergySearchCmd.Flags().String("clinical-status", "", "Clinical status (active, inactive, resolved)")
	allergySearchCmd.Flags().String("criticality", "", "Criticality (low, high, unable-to-assess)")
	allergySearchCmd.Flags().String("type", "", "Type (allergy, intolerance)")
	allergySearchCmd.Flags().String("category", "", "Category (food, medication, environment, biologic)")
	allergySearchCmd.Flags().Int("count", 20, "Maximum number of results")

	// List flags
	allergyListCmd.Flags().Int("count", 50, "Maximum number of results")
}
