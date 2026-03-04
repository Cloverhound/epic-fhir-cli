package cmd

import (
	"fmt"

	"github.com/Cloverhound/epic-fhir-cli/internal/config"
	"github.com/Cloverhound/epic-fhir-cli/internal/fhir"
	"github.com/Cloverhound/epic-fhir-cli/internal/output"
	"github.com/spf13/cobra"
)

var practitionerCmd = &cobra.Command{
	Use:     "practitioner",
	Aliases: []string{"provider"},
	Short:   "Practitioner resource operations",
	Long:    `Manage Practitioner resources - healthcare providers.`,
}

var practitionerGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get a practitioner by ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]

		client := fhir.NewClient(config.CurrentProfile)
		result, err := client.Get(fhir.ResourcePractitioner, id)
		if err != nil {
			return err
		}

		return output.Print(result, output.ParseFormat(GetOutput()))
	},
}

var practitionerSearchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for practitioners",
	Long: `Search for healthcare practitioners.

Examples:
  fhir-cli practitioner search --name "Smith"
  fhir-cli practitioner search --family "Smith" --given "John"
  fhir-cli practitioner search --identifier "NPI|1234567890"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		params := make(map[string]string)

		if name, _ := cmd.Flags().GetString("name"); name != "" {
			params["name"] = name
		}
		if family, _ := cmd.Flags().GetString("family"); family != "" {
			params["family"] = family
		}
		if given, _ := cmd.Flags().GetString("given"); given != "" {
			params["given"] = given
		}
		if identifier, _ := cmd.Flags().GetString("identifier"); identifier != "" {
			params["identifier"] = identifier
		}
		if active, _ := cmd.Flags().GetBool("active"); active {
			params["active"] = "true"
		}

		if count, _ := cmd.Flags().GetInt("count"); count > 0 {
			params["_count"] = fmt.Sprintf("%d", count)
		}

		client := fhir.NewClient(config.CurrentProfile)
		bundle, err := client.Search(fhir.ResourcePractitioner, params)
		if err != nil {
			return err
		}

		return output.Print(bundle, output.ParseFormat(GetOutput()))
	},
}

// PractitionerRole commands
var practitionerRoleCmd = &cobra.Command{
	Use:   "role",
	Short: "PractitionerRole operations",
	Long:  `Manage PractitionerRole resources - roles and specialties.`,
}

var practitionerRoleSearchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for practitioner roles",
	RunE: func(cmd *cobra.Command, args []string) error {
		params := make(map[string]string)

		if practitioner, _ := cmd.Flags().GetString("practitioner"); practitioner != "" {
			params["practitioner"] = practitioner
		}
		if specialty, _ := cmd.Flags().GetString("specialty"); specialty != "" {
			params["specialty"] = specialty
		}
		if organization, _ := cmd.Flags().GetString("organization"); organization != "" {
			params["organization"] = organization
		}

		if count, _ := cmd.Flags().GetInt("count"); count > 0 {
			params["_count"] = fmt.Sprintf("%d", count)
		}

		client := fhir.NewClient(config.CurrentProfile)
		bundle, err := client.Search(fhir.ResourcePractitionerRole, params)
		if err != nil {
			return err
		}

		return output.Print(bundle, output.ParseFormat(GetOutput()))
	},
}

// Organization commands
var organizationCmd = &cobra.Command{
	Use:     "organization",
	Aliases: []string{"org"},
	Short:   "Organization resource operations",
	Long:    `Manage Organization resources - healthcare organizations.`,
}

var organizationGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get an organization by ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]

		client := fhir.NewClient(config.CurrentProfile)
		result, err := client.Get(fhir.ResourceOrganization, id)
		if err != nil {
			return err
		}

		return output.Print(result, output.ParseFormat(GetOutput()))
	},
}

var organizationSearchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for organizations",
	RunE: func(cmd *cobra.Command, args []string) error {
		params := make(map[string]string)

		if name, _ := cmd.Flags().GetString("name"); name != "" {
			params["name"] = name
		}
		if identifier, _ := cmd.Flags().GetString("identifier"); identifier != "" {
			params["identifier"] = identifier
		}
		if orgType, _ := cmd.Flags().GetString("type"); orgType != "" {
			params["type"] = orgType
		}
		if address, _ := cmd.Flags().GetString("address"); address != "" {
			params["address"] = address
		}

		if count, _ := cmd.Flags().GetInt("count"); count > 0 {
			params["_count"] = fmt.Sprintf("%d", count)
		}

		client := fhir.NewClient(config.CurrentProfile)
		bundle, err := client.Search(fhir.ResourceOrganization, params)
		if err != nil {
			return err
		}

		return output.Print(bundle, output.ParseFormat(GetOutput()))
	},
}

// Location commands
var locationCmd = &cobra.Command{
	Use:   "location",
	Short: "Location resource operations",
	Long:  `Manage Location resources - physical places.`,
}

var locationGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get a location by ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]

		client := fhir.NewClient(config.CurrentProfile)
		result, err := client.Get(fhir.ResourceLocation, id)
		if err != nil {
			return err
		}

		return output.Print(result, output.ParseFormat(GetOutput()))
	},
}

var locationSearchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for locations",
	RunE: func(cmd *cobra.Command, args []string) error {
		params := make(map[string]string)

		if name, _ := cmd.Flags().GetString("name"); name != "" {
			params["name"] = name
		}
		if identifier, _ := cmd.Flags().GetString("identifier"); identifier != "" {
			params["identifier"] = identifier
		}
		if locType, _ := cmd.Flags().GetString("type"); locType != "" {
			params["type"] = locType
		}
		if address, _ := cmd.Flags().GetString("address"); address != "" {
			params["address"] = address
		}

		if count, _ := cmd.Flags().GetInt("count"); count > 0 {
			params["_count"] = fmt.Sprintf("%d", count)
		}

		client := fhir.NewClient(config.CurrentProfile)
		bundle, err := client.Search(fhir.ResourceLocation, params)
		if err != nil {
			return err
		}

		return output.Print(bundle, output.ParseFormat(GetOutput()))
	},
}

func init() {
	rootCmd.AddCommand(practitionerCmd)
	practitionerCmd.AddCommand(practitionerGetCmd)
	practitionerCmd.AddCommand(practitionerSearchCmd)
	practitionerCmd.AddCommand(practitionerRoleCmd)
	practitionerRoleCmd.AddCommand(practitionerRoleSearchCmd)

	rootCmd.AddCommand(organizationCmd)
	organizationCmd.AddCommand(organizationGetCmd)
	organizationCmd.AddCommand(organizationSearchCmd)

	rootCmd.AddCommand(locationCmd)
	locationCmd.AddCommand(locationGetCmd)
	locationCmd.AddCommand(locationSearchCmd)

	// Practitioner search flags
	practitionerSearchCmd.Flags().String("name", "", "Search by name")
	practitionerSearchCmd.Flags().String("family", "", "Search by family name")
	practitionerSearchCmd.Flags().String("given", "", "Search by given name")
	practitionerSearchCmd.Flags().String("identifier", "", "Search by identifier (system|value)")
	practitionerSearchCmd.Flags().Bool("active", false, "Only active practitioners")
	practitionerSearchCmd.Flags().Int("count", 20, "Maximum number of results")

	// PractitionerRole search flags
	practitionerRoleSearchCmd.Flags().String("practitioner", "", "Practitioner FHIR ID")
	practitionerRoleSearchCmd.Flags().String("specialty", "", "Specialty code")
	practitionerRoleSearchCmd.Flags().String("organization", "", "Organization FHIR ID")
	practitionerRoleSearchCmd.Flags().Int("count", 20, "Maximum number of results")

	// Organization search flags
	organizationSearchCmd.Flags().String("name", "", "Search by name")
	organizationSearchCmd.Flags().String("identifier", "", "Search by identifier")
	organizationSearchCmd.Flags().String("type", "", "Organization type")
	organizationSearchCmd.Flags().String("address", "", "Search by address")
	organizationSearchCmd.Flags().Int("count", 20, "Maximum number of results")

	// Location search flags
	locationSearchCmd.Flags().String("name", "", "Search by name")
	locationSearchCmd.Flags().String("identifier", "", "Search by identifier")
	locationSearchCmd.Flags().String("type", "", "Location type")
	locationSearchCmd.Flags().String("address", "", "Search by address")
	locationSearchCmd.Flags().Int("count", 20, "Maximum number of results")
}
