package cmd

import (
	"fmt"
	"sort"
	"strings"

	"github.com/jbogarin/fhir-cli/internal/config"
	"github.com/jbogarin/fhir-cli/internal/fhir"
	"github.com/jbogarin/fhir-cli/internal/output"
	"github.com/spf13/cobra"
)

var resourceCmd = &cobra.Command{
	Use:   "resource",
	Short: "Generic FHIR resource operations",
	Long: `Perform operations on any FHIR resource type.

Use this command for resources that don't have dedicated commands,
or when you need direct access to the FHIR API.`,
}

var resourceGetCmd = &cobra.Command{
	Use:   "get <resource-type> <id>",
	Short: "Get a resource by ID",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		resourceType := args[0]
		id := args[1]

		client := fhir.NewClient(config.CurrentProfile)
		result, err := client.Get(resourceType, id)
		if err != nil {
			return err
		}

		return output.Print(result, output.ParseFormat(GetOutput()))
	},
}

var resourceSearchCmd = &cobra.Command{
	Use:   "search <resource-type> [param=value...]",
	Short: "Search for resources",
	Long: `Search for FHIR resources with query parameters.

Parameters are specified as key=value pairs. Multiple parameters are ANDed together.

Examples:
  fhir-cli resource search Patient family=Smith
  fhir-cli resource search Observation patient=123 category=vital-signs
  fhir-cli resource search Condition patient=123 clinical-status=active`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		resourceType := args[0]
		params := make(map[string]string)

		// Parse remaining args as key=value pairs
		for _, arg := range args[1:] {
			parts := strings.SplitN(arg, "=", 2)
			if len(parts) == 2 {
				params[parts[0]] = parts[1]
			}
		}

		// Add count if specified
		if count, _ := cmd.Flags().GetInt("count"); count > 0 {
			params["_count"] = fmt.Sprintf("%d", count)
		}

		client := fhir.NewClient(config.CurrentProfile)
		bundle, err := client.Search(resourceType, params)
		if err != nil {
			return err
		}

		return output.Print(bundle, output.ParseFormat(GetOutput()))
	},
}

var resourceListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all supported resource types",
	RunE: func(cmd *cobra.Command, args []string) error {
		types := fhir.ListResourceTypes()
		sort.Strings(types)

		format := output.ParseFormat(GetOutput())

		if format == output.FormatTable {
			fmt.Println("Supported FHIR Resource Types:")
			fmt.Println("==============================")
			for _, t := range types {
				info, _ := fhir.GetResourceInfo(t)
				fmt.Printf("  %-28s %s\n", t, info.Description)
			}
			return nil
		}

		return output.Print(types, format)
	},
}

var resourceMetadataCmd = &cobra.Command{
	Use:   "metadata",
	Short: "Get server capability statement",
	Long:  `Retrieve the CapabilityStatement from the FHIR server.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := fhir.NewClient(config.CurrentProfile)
		result, err := client.GetMetadata()
		if err != nil {
			return err
		}

		return output.Print(result, output.ParseFormat(GetOutput()))
	},
}

var resourceOperationCmd = &cobra.Command{
	Use:   "operation <operation-name> [param=value...]",
	Short: "Execute a FHIR operation",
	Long: `Execute a FHIR operation ($operation) on the server.

Examples:
  fhir-cli resource operation everything patient=123
  fhir-cli resource operation validate-code system=http://loinc.org code=1234-5`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		operation := args[0]
		params := make(map[string]string)

		for _, arg := range args[1:] {
			parts := strings.SplitN(arg, "=", 2)
			if len(parts) == 2 {
				params[parts[0]] = parts[1]
			}
		}

		resType, _ := cmd.Flags().GetString("resource-type")
		id, _ := cmd.Flags().GetString("id")

		client := fhir.NewClient(config.CurrentProfile)
		result, err := client.Operation(resType, id, operation, params)
		if err != nil {
			return err
		}

		return output.Print(result, output.ParseFormat(GetOutput()))
	},
}

func init() {
	rootCmd.AddCommand(resourceCmd)
	resourceCmd.AddCommand(resourceGetCmd)
	resourceCmd.AddCommand(resourceSearchCmd)
	resourceCmd.AddCommand(resourceListCmd)
	resourceCmd.AddCommand(resourceMetadataCmd)
	resourceCmd.AddCommand(resourceOperationCmd)

	resourceSearchCmd.Flags().Int("count", 10, "Maximum number of results to return")
	resourceOperationCmd.Flags().String("resource-type", "", "Resource type for the operation")
	resourceOperationCmd.Flags().String("id", "", "Resource ID for instance-level operations")
}
