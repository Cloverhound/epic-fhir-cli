package cmd

import (
	"fmt"

	"github.com/jbogarin/fhir-cli/internal/config"
	"github.com/jbogarin/fhir-cli/internal/fhir"
	"github.com/jbogarin/fhir-cli/internal/output"
	"github.com/spf13/cobra"
)

var coverageCmd = &cobra.Command{
	Use:     "coverage",
	Aliases: []string{"insurance"},
	Short:   "Coverage resource operations",
	Long:    `Manage Coverage resources - insurance and payment information.`,
}

var coverageGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get a coverage record by ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]

		client := fhir.NewClient(config.CurrentProfile)
		result, err := client.Get(fhir.ResourceCoverage, id)
		if err != nil {
			return err
		}

		return output.Print(result, output.ParseFormat(GetOutput()))
	},
}

var coverageSearchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for coverage records",
	Long: `Search for insurance/coverage records.

Status values:
  - active, cancelled, draft, entered-in-error

Examples:
  fhir-cli coverage search --patient 123
  fhir-cli coverage search --patient 123 --status active`,
	RunE: func(cmd *cobra.Command, args []string) error {
		params := make(map[string]string)

		if patient, _ := cmd.Flags().GetString("patient"); patient != "" {
			params["patient"] = patient
		}
		if beneficiary, _ := cmd.Flags().GetString("beneficiary"); beneficiary != "" {
			params["beneficiary"] = beneficiary
		}
		if status, _ := cmd.Flags().GetString("status"); status != "" {
			params["status"] = status
		}
		if coverageType, _ := cmd.Flags().GetString("type"); coverageType != "" {
			params["type"] = coverageType
		}
		if payor, _ := cmd.Flags().GetString("payor"); payor != "" {
			params["payor"] = payor
		}

		if count, _ := cmd.Flags().GetInt("count"); count > 0 {
			params["_count"] = fmt.Sprintf("%d", count)
		}

		client := fhir.NewClient(config.CurrentProfile)
		bundle, err := client.Search(fhir.ResourceCoverage, params)
		if err != nil {
			return err
		}

		return output.Print(bundle, output.ParseFormat(GetOutput()))
	},
}

var coverageListCmd = &cobra.Command{
	Use:   "list <patient-id>",
	Short: "List coverage records for a patient",
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
		bundle, err := client.Search(fhir.ResourceCoverage, params)
		if err != nil {
			return err
		}

		return output.Print(bundle, output.ParseFormat(GetOutput()))
	},
}

// ExplanationOfBenefit commands
var eobCmd = &cobra.Command{
	Use:     "eob",
	Aliases: []string{"explanation-of-benefit"},
	Short:   "ExplanationOfBenefit resource operations",
	Long:    `Manage ExplanationOfBenefit resources - claim adjudication results.`,
}

var eobGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get an explanation of benefit by ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]

		client := fhir.NewClient(config.CurrentProfile)
		result, err := client.Get(fhir.ResourceExplanationOfBenefit, id)
		if err != nil {
			return err
		}

		return output.Print(result, output.ParseFormat(GetOutput()))
	},
}

var eobSearchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for explanations of benefit",
	RunE: func(cmd *cobra.Command, args []string) error {
		params := make(map[string]string)

		if patient, _ := cmd.Flags().GetString("patient"); patient != "" {
			params["patient"] = patient
		}
		if status, _ := cmd.Flags().GetString("status"); status != "" {
			params["status"] = status
		}
		if created, _ := cmd.Flags().GetString("created"); created != "" {
			params["created"] = created
		}

		if count, _ := cmd.Flags().GetInt("count"); count > 0 {
			params["_count"] = fmt.Sprintf("%d", count)
		}

		client := fhir.NewClient(config.CurrentProfile)
		bundle, err := client.Search(fhir.ResourceExplanationOfBenefit, params)
		if err != nil {
			return err
		}

		return output.Print(bundle, output.ParseFormat(GetOutput()))
	},
}

// ServiceRequest commands
var serviceRequestCmd = &cobra.Command{
	Use:     "servicerequest",
	Aliases: []string{"service-request", "order"},
	Short:   "ServiceRequest resource operations",
	Long:    `Manage ServiceRequest resources - orders and referrals.`,
}

var serviceRequestGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get a service request by ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]

		client := fhir.NewClient(config.CurrentProfile)
		result, err := client.Get(fhir.ResourceServiceRequest, id)
		if err != nil {
			return err
		}

		return output.Print(result, output.ParseFormat(GetOutput()))
	},
}

var serviceRequestSearchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for service requests",
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
		if category, _ := cmd.Flags().GetString("category"); category != "" {
			params["category"] = category
		}

		if count, _ := cmd.Flags().GetInt("count"); count > 0 {
			params["_count"] = fmt.Sprintf("%d", count)
		}

		client := fhir.NewClient(config.CurrentProfile)
		bundle, err := client.Search(fhir.ResourceServiceRequest, params)
		if err != nil {
			return err
		}

		return output.Print(bundle, output.ParseFormat(GetOutput()))
	},
}

func init() {
	rootCmd.AddCommand(coverageCmd)
	coverageCmd.AddCommand(coverageGetCmd)
	coverageCmd.AddCommand(coverageSearchCmd)
	coverageCmd.AddCommand(coverageListCmd)

	rootCmd.AddCommand(eobCmd)
	eobCmd.AddCommand(eobGetCmd)
	eobCmd.AddCommand(eobSearchCmd)

	rootCmd.AddCommand(serviceRequestCmd)
	serviceRequestCmd.AddCommand(serviceRequestGetCmd)
	serviceRequestCmd.AddCommand(serviceRequestSearchCmd)

	// Coverage search flags
	coverageSearchCmd.Flags().String("patient", "", "Patient FHIR ID")
	coverageSearchCmd.Flags().String("beneficiary", "", "Beneficiary FHIR ID")
	coverageSearchCmd.Flags().String("status", "", "Status (active, cancelled)")
	coverageSearchCmd.Flags().String("type", "", "Coverage type")
	coverageSearchCmd.Flags().String("payor", "", "Payor organization")
	coverageSearchCmd.Flags().Int("count", 20, "Maximum number of results")

	// Coverage list flags
	coverageListCmd.Flags().Int("count", 20, "Maximum number of results")

	// EOB search flags
	eobSearchCmd.Flags().String("patient", "", "Patient FHIR ID")
	eobSearchCmd.Flags().String("status", "", "Status")
	eobSearchCmd.Flags().String("created", "", "Created date range")
	eobSearchCmd.Flags().Int("count", 20, "Maximum number of results")

	// ServiceRequest search flags
	serviceRequestSearchCmd.Flags().String("patient", "", "Patient FHIR ID")
	serviceRequestSearchCmd.Flags().String("status", "", "Status (active, completed, cancelled)")
	serviceRequestSearchCmd.Flags().String("intent", "", "Intent (order, proposal, plan)")
	serviceRequestSearchCmd.Flags().String("category", "", "Category")
	serviceRequestSearchCmd.Flags().Int("count", 20, "Maximum number of results")
}
