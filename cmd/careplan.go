package cmd

import (
	"fmt"

	"github.com/Cloverhound/epic-fhir-cli/internal/config"
	"github.com/Cloverhound/epic-fhir-cli/internal/fhir"
	"github.com/Cloverhound/epic-fhir-cli/internal/output"
	"github.com/spf13/cobra"
)

var careplanCmd = &cobra.Command{
	Use:     "careplan",
	Aliases: []string{"care-plan"},
	Short:   "CarePlan resource operations",
	Long:    `Manage CarePlan resources - healthcare plans for patients.`,
}

var careplanGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get a care plan by ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]

		client := fhir.NewClient(config.CurrentProfile)
		result, err := client.Get(fhir.ResourceCarePlan, id)
		if err != nil {
			return err
		}

		return output.Print(result, output.ParseFormat(GetOutput()))
	},
}

var careplanSearchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for care plans",
	Long: `Search for care plan records.

Status values:
  - draft, active, on-hold, revoked, completed, entered-in-error, unknown

Intent values:
  - proposal, plan, order, option

Examples:
  fhir-cli careplan search --patient 123
  fhir-cli careplan search --patient 123 --status active`,
	RunE: func(cmd *cobra.Command, args []string) error {
		params := make(map[string]string)

		if patient, _ := cmd.Flags().GetString("patient"); patient != "" {
			params["patient"] = patient
		}
		if status, _ := cmd.Flags().GetString("status"); status != "" {
			params["status"] = status
		}
		if category, _ := cmd.Flags().GetString("category"); category != "" {
			params["category"] = category
		}
		if date, _ := cmd.Flags().GetString("date"); date != "" {
			params["date"] = date
		}

		if count, _ := cmd.Flags().GetInt("count"); count > 0 {
			params["_count"] = fmt.Sprintf("%d", count)
		}

		client := fhir.NewClient(config.CurrentProfile)
		bundle, err := client.Search(fhir.ResourceCarePlan, params)
		if err != nil {
			return err
		}

		return output.Print(bundle, output.ParseFormat(GetOutput()))
	},
}

var careplanListCmd = &cobra.Command{
	Use:   "list <patient-id>",
	Short: "List care plans for a patient",
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
		bundle, err := client.Search(fhir.ResourceCarePlan, params)
		if err != nil {
			return err
		}

		return output.Print(bundle, output.ParseFormat(GetOutput()))
	},
}

// CareTeam commands
var careteamCmd = &cobra.Command{
	Use:     "careteam",
	Aliases: []string{"care-team"},
	Short:   "CareTeam resource operations",
	Long:    `Manage CareTeam resources - team of practitioners caring for a patient.`,
}

var careteamSearchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for care teams",
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
		bundle, err := client.Search(fhir.ResourceCareTeam, params)
		if err != nil {
			return err
		}

		return output.Print(bundle, output.ParseFormat(GetOutput()))
	},
}

// Goal commands
var goalCmd = &cobra.Command{
	Use:   "goal",
	Short: "Goal resource operations",
	Long:  `Manage Goal resources - desired health objectives for patients.`,
}

var goalSearchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for goals",
	RunE: func(cmd *cobra.Command, args []string) error {
		params := make(map[string]string)

		if patient, _ := cmd.Flags().GetString("patient"); patient != "" {
			params["patient"] = patient
		}
		if status, _ := cmd.Flags().GetString("lifecycle-status"); status != "" {
			params["lifecycle-status"] = status
		}

		if count, _ := cmd.Flags().GetInt("count"); count > 0 {
			params["_count"] = fmt.Sprintf("%d", count)
		}

		client := fhir.NewClient(config.CurrentProfile)
		bundle, err := client.Search(fhir.ResourceGoal, params)
		if err != nil {
			return err
		}

		return output.Print(bundle, output.ParseFormat(GetOutput()))
	},
}

func init() {
	rootCmd.AddCommand(careplanCmd)
	careplanCmd.AddCommand(careplanGetCmd)
	careplanCmd.AddCommand(careplanSearchCmd)
	careplanCmd.AddCommand(careplanListCmd)

	rootCmd.AddCommand(careteamCmd)
	careteamCmd.AddCommand(careteamSearchCmd)

	rootCmd.AddCommand(goalCmd)
	goalCmd.AddCommand(goalSearchCmd)

	// CarePlan search flags
	careplanSearchCmd.Flags().String("patient", "", "Patient FHIR ID")
	careplanSearchCmd.Flags().String("status", "", "Status (active, completed, draft)")
	careplanSearchCmd.Flags().String("category", "", "Category")
	careplanSearchCmd.Flags().String("date", "", "Date range")
	careplanSearchCmd.Flags().Int("count", 20, "Maximum number of results")

	// CarePlan list flags
	careplanListCmd.Flags().Int("count", 20, "Maximum number of results")

	// CareTeam search flags
	careteamSearchCmd.Flags().String("patient", "", "Patient FHIR ID")
	careteamSearchCmd.Flags().String("status", "", "Status")
	careteamSearchCmd.Flags().Int("count", 20, "Maximum number of results")

	// Goal search flags
	goalSearchCmd.Flags().String("patient", "", "Patient FHIR ID")
	goalSearchCmd.Flags().String("lifecycle-status", "", "Lifecycle status")
	goalSearchCmd.Flags().Int("count", 20, "Maximum number of results")
}
