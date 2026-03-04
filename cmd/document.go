package cmd

import (
	"fmt"

	"github.com/Cloverhound/epic-fhir-cli/internal/config"
	"github.com/Cloverhound/epic-fhir-cli/internal/fhir"
	"github.com/Cloverhound/epic-fhir-cli/internal/output"
	"github.com/spf13/cobra"
)

var documentCmd = &cobra.Command{
	Use:     "document",
	Aliases: []string{"doc"},
	Short:   "DocumentReference resource operations",
	Long:    `Manage DocumentReference resources - clinical documents.`,
}

var documentGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get a document reference by ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]

		client := fhir.NewClient(config.CurrentProfile)
		result, err := client.Get(fhir.ResourceDocumentReference, id)
		if err != nil {
			return err
		}

		return output.Print(result, output.ParseFormat(GetOutput()))
	},
}

var documentSearchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for document references",
	Long: `Search for clinical documents.

Status values:
  - current: The document is current
  - superseded: The document has been superseded
  - entered-in-error: The document was entered in error

Type examples (LOINC codes):
  - 34133-9: Summary of episode note
  - 18842-5: Discharge summary
  - 11502-2: Laboratory report

Examples:
  fhir-cli document search --patient 123
  fhir-cli document search --patient 123 --type 34133-9
  fhir-cli document search --patient 123 --category clinical-note`,
	RunE: func(cmd *cobra.Command, args []string) error {
		params := make(map[string]string)

		if patient, _ := cmd.Flags().GetString("patient"); patient != "" {
			params["patient"] = patient
		}
		if status, _ := cmd.Flags().GetString("status"); status != "" {
			params["status"] = status
		}
		if docType, _ := cmd.Flags().GetString("type"); docType != "" {
			params["type"] = docType
		}
		if category, _ := cmd.Flags().GetString("category"); category != "" {
			params["category"] = category
		}
		if date, _ := cmd.Flags().GetString("date"); date != "" {
			params["date"] = date
		}
		if period, _ := cmd.Flags().GetString("period"); period != "" {
			params["period"] = period
		}

		if count, _ := cmd.Flags().GetInt("count"); count > 0 {
			params["_count"] = fmt.Sprintf("%d", count)
		}

		client := fhir.NewClient(config.CurrentProfile)
		bundle, err := client.Search(fhir.ResourceDocumentReference, params)
		if err != nil {
			return err
		}

		return output.Print(bundle, output.ParseFormat(GetOutput()))
	},
}

var documentListCmd = &cobra.Command{
	Use:   "list <patient-id>",
	Short: "List documents for a patient",
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
		bundle, err := client.Search(fhir.ResourceDocumentReference, params)
		if err != nil {
			return err
		}

		return output.Print(bundle, output.ParseFormat(GetOutput()))
	},
}

// Binary command for fetching document content
var binaryCmd = &cobra.Command{
	Use:   "binary",
	Short: "Binary resource operations",
	Long:  `Manage Binary resources - raw document content.`,
}

var binaryGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get binary content by ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]

		client := fhir.NewClient(config.CurrentProfile)
		result, err := client.Get(fhir.ResourceBinary, id)
		if err != nil {
			return err
		}

		return output.Print(result, output.ParseFormat(GetOutput()))
	},
}

func init() {
	rootCmd.AddCommand(documentCmd)
	documentCmd.AddCommand(documentGetCmd)
	documentCmd.AddCommand(documentSearchCmd)
	documentCmd.AddCommand(documentListCmd)

	rootCmd.AddCommand(binaryCmd)
	binaryCmd.AddCommand(binaryGetCmd)

	// Search flags
	documentSearchCmd.Flags().String("patient", "", "Patient FHIR ID")
	documentSearchCmd.Flags().String("status", "", "Status (current, superseded, entered-in-error)")
	documentSearchCmd.Flags().String("type", "", "Document type (LOINC code)")
	documentSearchCmd.Flags().String("category", "", "Document category")
	documentSearchCmd.Flags().String("date", "", "Document date range")
	documentSearchCmd.Flags().String("period", "", "Period covered by document")
	documentSearchCmd.Flags().Int("count", 20, "Maximum number of results")

	// List flags
	documentListCmd.Flags().Int("count", 20, "Maximum number of results")
}
