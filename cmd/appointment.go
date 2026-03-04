package cmd

import (
	"fmt"

	"github.com/Cloverhound/epic-fhir-cli/internal/config"
	"github.com/Cloverhound/epic-fhir-cli/internal/fhir"
	"github.com/Cloverhound/epic-fhir-cli/internal/output"
	"github.com/spf13/cobra"
)

var appointmentCmd = &cobra.Command{
	Use:     "appointment",
	Aliases: []string{"appt"},
	Short:   "Appointment resource operations",
	Long:    `Manage Appointment resources - scheduled healthcare events.`,
}

var appointmentGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get an appointment by ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]

		client := fhir.NewClient(config.CurrentProfile)
		result, err := client.Get(fhir.ResourceAppointment, id)
		if err != nil {
			return err
		}

		return output.Print(result, output.ParseFormat(GetOutput()))
	},
}

var appointmentSearchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for appointments",
	Long: `Search for appointment records.

Status values:
  - proposed, pending, booked, arrived, fulfilled, cancelled, noshow, entered-in-error, checked-in, waitlist

Examples:
  fhir-cli appointment search --patient 123
  fhir-cli appointment search --patient 123 --status booked
  fhir-cli appointment search --patient 123 --date ge2024-01-01`,
	RunE: func(cmd *cobra.Command, args []string) error {
		params := make(map[string]string)

		if patient, _ := cmd.Flags().GetString("patient"); patient != "" {
			params["patient"] = patient
		}
		if status, _ := cmd.Flags().GetString("status"); status != "" {
			params["status"] = status
		}
		if date, _ := cmd.Flags().GetString("date"); date != "" {
			params["date"] = date
		}
		if practitioner, _ := cmd.Flags().GetString("practitioner"); practitioner != "" {
			params["practitioner"] = practitioner
		}
		if location, _ := cmd.Flags().GetString("location"); location != "" {
			params["location"] = location
		}
		if serviceType, _ := cmd.Flags().GetString("service-type"); serviceType != "" {
			params["service-type"] = serviceType
		}

		if count, _ := cmd.Flags().GetInt("count"); count > 0 {
			params["_count"] = fmt.Sprintf("%d", count)
		}

		client := fhir.NewClient(config.CurrentProfile)
		bundle, err := client.Search(fhir.ResourceAppointment, params)
		if err != nil {
			return err
		}

		return output.Print(bundle, output.ParseFormat(GetOutput()))
	},
}

var appointmentListCmd = &cobra.Command{
	Use:   "list <patient-id>",
	Short: "List appointments for a patient",
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
		bundle, err := client.Search(fhir.ResourceAppointment, params)
		if err != nil {
			return err
		}

		return output.Print(bundle, output.ParseFormat(GetOutput()))
	},
}

var appointmentUpcomingCmd = &cobra.Command{
	Use:   "upcoming <patient-id>",
	Short: "Get upcoming appointments for a patient",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		patientID := args[0]
		params := map[string]string{
			"patient": patientID,
			"status":  "booked",
			"date":    "ge" + "2024-01-01", // This should be dynamic
		}

		if count, _ := cmd.Flags().GetInt("count"); count > 0 {
			params["_count"] = fmt.Sprintf("%d", count)
		}

		client := fhir.NewClient(config.CurrentProfile)
		bundle, err := client.Search(fhir.ResourceAppointment, params)
		if err != nil {
			return err
		}

		return output.Print(bundle, output.ParseFormat(GetOutput()))
	},
}

// Schedule commands
var scheduleCmd = &cobra.Command{
	Use:   "schedule",
	Short: "Schedule resource operations",
	Long:  `Manage Schedule resources - provider availability.`,
}

var scheduleGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get a schedule by ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]

		client := fhir.NewClient(config.CurrentProfile)
		result, err := client.Get(fhir.ResourceSchedule, id)
		if err != nil {
			return err
		}

		return output.Print(result, output.ParseFormat(GetOutput()))
	},
}

var scheduleSearchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for schedules",
	RunE: func(cmd *cobra.Command, args []string) error {
		params := make(map[string]string)

		if actor, _ := cmd.Flags().GetString("actor"); actor != "" {
			params["actor"] = actor
		}
		if date, _ := cmd.Flags().GetString("date"); date != "" {
			params["date"] = date
		}

		if count, _ := cmd.Flags().GetInt("count"); count > 0 {
			params["_count"] = fmt.Sprintf("%d", count)
		}

		client := fhir.NewClient(config.CurrentProfile)
		bundle, err := client.Search(fhir.ResourceSchedule, params)
		if err != nil {
			return err
		}

		return output.Print(bundle, output.ParseFormat(GetOutput()))
	},
}

// Slot commands
var slotCmd = &cobra.Command{
	Use:   "slot",
	Short: "Slot resource operations",
	Long:  `Manage Slot resources - available appointment slots.`,
}

var slotSearchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for available slots",
	RunE: func(cmd *cobra.Command, args []string) error {
		params := make(map[string]string)

		if schedule, _ := cmd.Flags().GetString("schedule"); schedule != "" {
			params["schedule"] = schedule
		}
		if status, _ := cmd.Flags().GetString("status"); status != "" {
			params["status"] = status
		}
		if start, _ := cmd.Flags().GetString("start"); start != "" {
			params["start"] = start
		}

		if count, _ := cmd.Flags().GetInt("count"); count > 0 {
			params["_count"] = fmt.Sprintf("%d", count)
		}

		client := fhir.NewClient(config.CurrentProfile)
		bundle, err := client.Search(fhir.ResourceSlot, params)
		if err != nil {
			return err
		}

		return output.Print(bundle, output.ParseFormat(GetOutput()))
	},
}

func init() {
	rootCmd.AddCommand(appointmentCmd)
	appointmentCmd.AddCommand(appointmentGetCmd)
	appointmentCmd.AddCommand(appointmentSearchCmd)
	appointmentCmd.AddCommand(appointmentListCmd)
	appointmentCmd.AddCommand(appointmentUpcomingCmd)

	rootCmd.AddCommand(scheduleCmd)
	scheduleCmd.AddCommand(scheduleGetCmd)
	scheduleCmd.AddCommand(scheduleSearchCmd)

	rootCmd.AddCommand(slotCmd)
	slotCmd.AddCommand(slotSearchCmd)

	// Appointment search flags
	appointmentSearchCmd.Flags().String("patient", "", "Patient FHIR ID")
	appointmentSearchCmd.Flags().String("status", "", "Status (booked, pending, cancelled)")
	appointmentSearchCmd.Flags().String("date", "", "Date range")
	appointmentSearchCmd.Flags().String("practitioner", "", "Practitioner FHIR ID")
	appointmentSearchCmd.Flags().String("location", "", "Location FHIR ID")
	appointmentSearchCmd.Flags().String("service-type", "", "Service type")
	appointmentSearchCmd.Flags().Int("count", 20, "Maximum number of results")

	// Appointment list flags
	appointmentListCmd.Flags().Int("count", 20, "Maximum number of results")

	// Appointment upcoming flags
	appointmentUpcomingCmd.Flags().Int("count", 10, "Maximum number of results")

	// Schedule search flags
	scheduleSearchCmd.Flags().String("actor", "", "Actor (practitioner/location) FHIR ID")
	scheduleSearchCmd.Flags().String("date", "", "Date range")
	scheduleSearchCmd.Flags().Int("count", 20, "Maximum number of results")

	// Slot search flags
	slotSearchCmd.Flags().String("schedule", "", "Schedule FHIR ID")
	slotSearchCmd.Flags().String("status", "", "Status (free, busy)")
	slotSearchCmd.Flags().String("start", "", "Start time")
	slotSearchCmd.Flags().Int("count", 50, "Maximum number of results")
}
