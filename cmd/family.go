package cmd

import (
	"fmt"

	"github.com/Cloverhound/epic-fhir-cli/internal/config"
	"github.com/Cloverhound/epic-fhir-cli/internal/fhir"
	"github.com/Cloverhound/epic-fhir-cli/internal/output"
	"github.com/spf13/cobra"
)

var familyHistoryCmd = &cobra.Command{
	Use:     "familyhistory",
	Aliases: []string{"family-history", "fmh"},
	Short:   "FamilyMemberHistory resource operations",
	Long:    `Manage FamilyMemberHistory resources - family health history.`,
}

var familyHistoryGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get a family member history by ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]

		client := fhir.NewClient(config.CurrentProfile)
		result, err := client.Get(fhir.ResourceFamilyMemberHistory, id)
		if err != nil {
			return err
		}

		return output.Print(result, output.ParseFormat(GetOutput()))
	},
}

var familyHistorySearchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for family member history records",
	Long: `Search for family member health history records.

Status values:
  - partial, completed, entered-in-error, health-unknown

Examples:
  fhir-cli familyhistory search --patient 123
  fhir-cli familyhistory search --patient 123 --relationship parent`,
	RunE: func(cmd *cobra.Command, args []string) error {
		params := make(map[string]string)

		if patient, _ := cmd.Flags().GetString("patient"); patient != "" {
			params["patient"] = patient
		}
		if status, _ := cmd.Flags().GetString("status"); status != "" {
			params["status"] = status
		}
		if relationship, _ := cmd.Flags().GetString("relationship"); relationship != "" {
			params["relationship"] = relationship
		}
		if code, _ := cmd.Flags().GetString("code"); code != "" {
			params["code"] = code
		}

		if count, _ := cmd.Flags().GetInt("count"); count > 0 {
			params["_count"] = fmt.Sprintf("%d", count)
		}

		client := fhir.NewClient(config.CurrentProfile)
		bundle, err := client.Search(fhir.ResourceFamilyMemberHistory, params)
		if err != nil {
			return err
		}

		return output.Print(bundle, output.ParseFormat(GetOutput()))
	},
}

var familyHistoryListCmd = &cobra.Command{
	Use:   "list <patient-id>",
	Short: "List family member history for a patient",
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
		bundle, err := client.Search(fhir.ResourceFamilyMemberHistory, params)
		if err != nil {
			return err
		}

		return output.Print(bundle, output.ParseFormat(GetOutput()))
	},
}

// RelatedPerson commands
var relatedPersonCmd = &cobra.Command{
	Use:     "relatedperson",
	Aliases: []string{"related-person", "contact"},
	Short:   "RelatedPerson resource operations",
	Long:    `Manage RelatedPerson resources - patient contacts and relationships.`,
}

var relatedPersonGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get a related person by ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]

		client := fhir.NewClient(config.CurrentProfile)
		result, err := client.Get(fhir.ResourceRelatedPerson, id)
		if err != nil {
			return err
		}

		return output.Print(result, output.ParseFormat(GetOutput()))
	},
}

var relatedPersonSearchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for related persons",
	RunE: func(cmd *cobra.Command, args []string) error {
		params := make(map[string]string)

		if patient, _ := cmd.Flags().GetString("patient"); patient != "" {
			params["patient"] = patient
		}
		if name, _ := cmd.Flags().GetString("name"); name != "" {
			params["name"] = name
		}
		if relationship, _ := cmd.Flags().GetString("relationship"); relationship != "" {
			params["relationship"] = relationship
		}

		if count, _ := cmd.Flags().GetInt("count"); count > 0 {
			params["_count"] = fmt.Sprintf("%d", count)
		}

		client := fhir.NewClient(config.CurrentProfile)
		bundle, err := client.Search(fhir.ResourceRelatedPerson, params)
		if err != nil {
			return err
		}

		return output.Print(bundle, output.ParseFormat(GetOutput()))
	},
}

// Consent commands
var consentCmd = &cobra.Command{
	Use:   "consent",
	Short: "Consent resource operations",
	Long:  `Manage Consent resources - patient privacy preferences.`,
}

var consentGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get a consent record by ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]

		client := fhir.NewClient(config.CurrentProfile)
		result, err := client.Get(fhir.ResourceConsent, id)
		if err != nil {
			return err
		}

		return output.Print(result, output.ParseFormat(GetOutput()))
	},
}

var consentSearchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for consent records",
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

		if count, _ := cmd.Flags().GetInt("count"); count > 0 {
			params["_count"] = fmt.Sprintf("%d", count)
		}

		client := fhir.NewClient(config.CurrentProfile)
		bundle, err := client.Search(fhir.ResourceConsent, params)
		if err != nil {
			return err
		}

		return output.Print(bundle, output.ParseFormat(GetOutput()))
	},
}

// Provenance commands
var provenanceCmd = &cobra.Command{
	Use:   "provenance",
	Short: "Provenance resource operations",
	Long:  `Manage Provenance resources - audit trail and data lineage.`,
}

var provenanceSearchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for provenance records",
	RunE: func(cmd *cobra.Command, args []string) error {
		params := make(map[string]string)

		if target, _ := cmd.Flags().GetString("target"); target != "" {
			params["target"] = target
		}
		if patient, _ := cmd.Flags().GetString("patient"); patient != "" {
			params["patient"] = patient
		}
		if recorded, _ := cmd.Flags().GetString("recorded"); recorded != "" {
			params["recorded"] = recorded
		}

		if count, _ := cmd.Flags().GetInt("count"); count > 0 {
			params["_count"] = fmt.Sprintf("%d", count)
		}

		client := fhir.NewClient(config.CurrentProfile)
		bundle, err := client.Search(fhir.ResourceProvenance, params)
		if err != nil {
			return err
		}

		return output.Print(bundle, output.ParseFormat(GetOutput()))
	},
}

// QuestionnaireResponse commands
var questionnaireResponseCmd = &cobra.Command{
	Use:     "questionnaireresponse",
	Aliases: []string{"questionnaire-response", "qr"},
	Short:   "QuestionnaireResponse resource operations",
	Long:    `Manage QuestionnaireResponse resources - completed questionnaire answers.`,
}

var questionnaireResponseGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get a questionnaire response by ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]

		client := fhir.NewClient(config.CurrentProfile)
		result, err := client.Get(fhir.ResourceQuestionnaireResponse, id)
		if err != nil {
			return err
		}

		return output.Print(result, output.ParseFormat(GetOutput()))
	},
}

var questionnaireResponseSearchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for questionnaire responses",
	RunE: func(cmd *cobra.Command, args []string) error {
		params := make(map[string]string)

		if patient, _ := cmd.Flags().GetString("patient"); patient != "" {
			params["patient"] = patient
		}
		if status, _ := cmd.Flags().GetString("status"); status != "" {
			params["status"] = status
		}
		if authored, _ := cmd.Flags().GetString("authored"); authored != "" {
			params["authored"] = authored
		}

		if count, _ := cmd.Flags().GetInt("count"); count > 0 {
			params["_count"] = fmt.Sprintf("%d", count)
		}

		client := fhir.NewClient(config.CurrentProfile)
		bundle, err := client.Search(fhir.ResourceQuestionnaireResponse, params)
		if err != nil {
			return err
		}

		return output.Print(bundle, output.ParseFormat(GetOutput()))
	},
}

func init() {
	rootCmd.AddCommand(familyHistoryCmd)
	familyHistoryCmd.AddCommand(familyHistoryGetCmd)
	familyHistoryCmd.AddCommand(familyHistorySearchCmd)
	familyHistoryCmd.AddCommand(familyHistoryListCmd)

	rootCmd.AddCommand(relatedPersonCmd)
	relatedPersonCmd.AddCommand(relatedPersonGetCmd)
	relatedPersonCmd.AddCommand(relatedPersonSearchCmd)

	rootCmd.AddCommand(consentCmd)
	consentCmd.AddCommand(consentGetCmd)
	consentCmd.AddCommand(consentSearchCmd)

	rootCmd.AddCommand(provenanceCmd)
	provenanceCmd.AddCommand(provenanceSearchCmd)

	rootCmd.AddCommand(questionnaireResponseCmd)
	questionnaireResponseCmd.AddCommand(questionnaireResponseGetCmd)
	questionnaireResponseCmd.AddCommand(questionnaireResponseSearchCmd)

	// FamilyMemberHistory flags
	familyHistorySearchCmd.Flags().String("patient", "", "Patient FHIR ID")
	familyHistorySearchCmd.Flags().String("status", "", "Status (completed, partial)")
	familyHistorySearchCmd.Flags().String("relationship", "", "Relationship to patient")
	familyHistorySearchCmd.Flags().String("code", "", "Condition code")
	familyHistorySearchCmd.Flags().Int("count", 20, "Maximum number of results")
	familyHistoryListCmd.Flags().Int("count", 50, "Maximum number of results")

	// RelatedPerson flags
	relatedPersonSearchCmd.Flags().String("patient", "", "Patient FHIR ID")
	relatedPersonSearchCmd.Flags().String("name", "", "Search by name")
	relatedPersonSearchCmd.Flags().String("relationship", "", "Relationship type")
	relatedPersonSearchCmd.Flags().Int("count", 20, "Maximum number of results")

	// Consent flags
	consentSearchCmd.Flags().String("patient", "", "Patient FHIR ID")
	consentSearchCmd.Flags().String("status", "", "Status (active, inactive, draft)")
	consentSearchCmd.Flags().String("category", "", "Category")
	consentSearchCmd.Flags().Int("count", 20, "Maximum number of results")

	// Provenance flags
	provenanceSearchCmd.Flags().String("target", "", "Target resource reference")
	provenanceSearchCmd.Flags().String("patient", "", "Patient FHIR ID")
	provenanceSearchCmd.Flags().String("recorded", "", "Recorded date range")
	provenanceSearchCmd.Flags().Int("count", 20, "Maximum number of results")

	// QuestionnaireResponse flags
	questionnaireResponseSearchCmd.Flags().String("patient", "", "Patient FHIR ID")
	questionnaireResponseSearchCmd.Flags().String("status", "", "Status (completed, in-progress)")
	questionnaireResponseSearchCmd.Flags().String("authored", "", "Authored date range")
	questionnaireResponseSearchCmd.Flags().Int("count", 20, "Maximum number of results")
}
