package cmd

import (
	"fmt"

	"github.com/jbogarin/fhir-cli/internal/config"
	"github.com/jbogarin/fhir-cli/internal/fhir"
	"github.com/jbogarin/fhir-cli/internal/output"
	"github.com/spf13/cobra"
)

var patientCmd = &cobra.Command{
	Use:   "patient",
	Short: "Patient resource operations",
	Long:  `Manage Patient resources - demographics and administrative information.`,
}

var patientGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get a patient by FHIR ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]

		client := fhir.NewClient(config.CurrentProfile)
		result, err := client.Get(fhir.ResourcePatient, id)
		if err != nil {
			return err
		}

		return output.Print(result, output.ParseFormat(GetOutput()))
	},
}

var patientSearchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for patients",
	Long: `Search for patients using various criteria.

Examples:
  fhir-cli patient search --name "John Smith"
  fhir-cli patient search --family Smith --given John
  fhir-cli patient search --birthdate 1980-01-15
  fhir-cli patient search --identifier "MRN|12345"
  fhir-cli patient search --gender male --address-city "New York"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		params := make(map[string]string)

		// Collect all search parameters
		if name, _ := cmd.Flags().GetString("name"); name != "" {
			params["name"] = name
		}
		if family, _ := cmd.Flags().GetString("family"); family != "" {
			params["family"] = family
		}
		if given, _ := cmd.Flags().GetString("given"); given != "" {
			params["given"] = given
		}
		if birthdate, _ := cmd.Flags().GetString("birthdate"); birthdate != "" {
			params["birthdate"] = birthdate
		}
		if gender, _ := cmd.Flags().GetString("gender"); gender != "" {
			params["gender"] = gender
		}
		if identifier, _ := cmd.Flags().GetString("identifier"); identifier != "" {
			params["identifier"] = identifier
		}
		if address, _ := cmd.Flags().GetString("address"); address != "" {
			params["address"] = address
		}
		if city, _ := cmd.Flags().GetString("address-city"); city != "" {
			params["address-city"] = city
		}
		if state, _ := cmd.Flags().GetString("address-state"); state != "" {
			params["address-state"] = state
		}
		if postalCode, _ := cmd.Flags().GetString("address-postalcode"); postalCode != "" {
			params["address-postalcode"] = postalCode
		}
		if phone, _ := cmd.Flags().GetString("phone"); phone != "" {
			params["phone"] = phone
		}
		if email, _ := cmd.Flags().GetString("email"); email != "" {
			params["email"] = email
		}

		if count, _ := cmd.Flags().GetInt("count"); count > 0 {
			params["_count"] = fmt.Sprintf("%d", count)
		}

		if len(params) == 0 || (len(params) == 1 && params["_count"] != "") {
			return fmt.Errorf("at least one search parameter is required")
		}

		client := fhir.NewClient(config.CurrentProfile)
		bundle, err := client.Search(fhir.ResourcePatient, params)
		if err != nil {
			return err
		}

		return output.Print(bundle, output.ParseFormat(GetOutput()))
	},
}

var patientEverythingCmd = &cobra.Command{
	Use:   "everything <id>",
	Short: "Get all data for a patient ($everything operation)",
	Long: `Retrieve all resources related to a patient using the $everything operation.

This returns a Bundle containing all resources that reference the patient.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]

		client := fhir.NewClient(config.CurrentProfile)
		result, err := client.Operation(fhir.ResourcePatient, id, "everything", nil)
		if err != nil {
			return err
		}

		return output.Print(result, output.ParseFormat(GetOutput()))
	},
}

func init() {
	rootCmd.AddCommand(patientCmd)
	patientCmd.AddCommand(patientGetCmd)
	patientCmd.AddCommand(patientSearchCmd)
	patientCmd.AddCommand(patientEverythingCmd)

	// Search flags
	patientSearchCmd.Flags().String("name", "", "Search by full name")
	patientSearchCmd.Flags().String("family", "", "Search by family/last name")
	patientSearchCmd.Flags().String("given", "", "Search by given/first name")
	patientSearchCmd.Flags().String("birthdate", "", "Search by birth date (YYYY-MM-DD)")
	patientSearchCmd.Flags().String("gender", "", "Search by gender (male, female, other, unknown)")
	patientSearchCmd.Flags().String("identifier", "", "Search by identifier (system|value)")
	patientSearchCmd.Flags().String("address", "", "Search by address text")
	patientSearchCmd.Flags().String("address-city", "", "Search by city")
	patientSearchCmd.Flags().String("address-state", "", "Search by state")
	patientSearchCmd.Flags().String("address-postalcode", "", "Search by postal code")
	patientSearchCmd.Flags().String("phone", "", "Search by phone number")
	patientSearchCmd.Flags().String("email", "", "Search by email")
	patientSearchCmd.Flags().Int("count", 10, "Maximum number of results")
}
