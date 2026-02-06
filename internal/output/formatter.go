package output

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
	"gopkg.in/yaml.v3"
)

type Format string

const (
	FormatJSON  Format = "json"
	FormatTable Format = "table"
	FormatYAML  Format = "yaml"
)

func ParseFormat(s string) Format {
	switch strings.ToLower(s) {
	case "table":
		return FormatTable
	case "yaml", "yml":
		return FormatYAML
	default:
		return FormatJSON
	}
}

func Print(data interface{}, format Format) error {
	switch format {
	case FormatJSON:
		return printJSON(data)
	case FormatYAML:
		return printYAML(data)
	case FormatTable:
		return printTable(data)
	default:
		return printJSON(data)
	}
}

func printJSON(data interface{}) error {
	var output []byte
	var err error

	switch v := data.(type) {
	case []byte:
		// Pretty print if it's already JSON bytes
		var parsed interface{}
		if json.Unmarshal(v, &parsed) == nil {
			output, err = json.MarshalIndent(parsed, "", "  ")
		} else {
			output = v
		}
	case json.RawMessage:
		var parsed interface{}
		if json.Unmarshal(v, &parsed) == nil {
			output, err = json.MarshalIndent(parsed, "", "  ")
		} else {
			output = v
		}
	default:
		output, err = json.MarshalIndent(data, "", "  ")
	}

	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	fmt.Println(string(output))
	return nil
}

func printYAML(data interface{}) error {
	var toMarshal interface{}

	switch v := data.(type) {
	case []byte:
		if json.Unmarshal(v, &toMarshal) != nil {
			toMarshal = string(v)
		}
	case json.RawMessage:
		if json.Unmarshal(v, &toMarshal) != nil {
			toMarshal = string(v)
		}
	default:
		toMarshal = data
	}

	output, err := yaml.Marshal(toMarshal)
	if err != nil {
		return fmt.Errorf("failed to marshal YAML: %w", err)
	}

	fmt.Println(string(output))
	return nil
}

func printTable(data interface{}) error {
	// Convert data to a map for processing
	var parsed map[string]interface{}

	switch v := data.(type) {
	case []byte:
		if err := json.Unmarshal(v, &parsed); err != nil {
			return printJSON(data) // Fall back to JSON if not a valid object
		}
	case json.RawMessage:
		if err := json.Unmarshal(v, &parsed); err != nil {
			return printJSON(data)
		}
	default:
		jsonBytes, err := json.Marshal(data)
		if err != nil {
			return printJSON(data)
		}
		if err := json.Unmarshal(jsonBytes, &parsed); err != nil {
			return printJSON(data)
		}
	}

	// Check if it's a Bundle
	if resourceType, ok := parsed["resourceType"].(string); ok {
		switch resourceType {
		case "Bundle":
			return printBundleTable(parsed)
		case "Patient":
			return printPatientTable(parsed)
		case "Observation":
			return printObservationTable(parsed)
		case "Condition":
			return printConditionTable(parsed)
		case "MedicationRequest":
			return printMedicationRequestTable(parsed)
		case "Encounter":
			return printEncounterTable(parsed)
		default:
			return printGenericResourceTable(parsed)
		}
	}

	return printJSON(data)
}

func printBundleTable(bundle map[string]interface{}) error {
	entries, ok := bundle["entry"].([]interface{})
	if !ok || len(entries) == 0 {
		fmt.Println("No entries found")
		return nil
	}

	// Get total if available
	if total, ok := bundle["total"].(float64); ok {
		fmt.Printf("Total: %d\n\n", int(total))
	}

	// Determine resource type from first entry
	firstEntry, ok := entries[0].(map[string]interface{})
	if !ok {
		return printJSON(bundle)
	}

	resource, ok := firstEntry["resource"].(map[string]interface{})
	if !ok {
		return printJSON(bundle)
	}

	resourceType, _ := resource["resourceType"].(string)

	// Create table based on resource type
	table := tablewriter.NewWriter(os.Stdout)

	switch resourceType {
	case "Patient":
		return printPatientBundleTable(entries)
	case "Observation":
		return printObservationBundleTable(entries)
	case "Condition":
		return printConditionBundleTable(entries)
	case "MedicationRequest":
		return printMedicationRequestBundleTable(entries)
	case "Encounter":
		return printEncounterBundleTable(entries)
	default:
		table.SetHeader([]string{"ID", "Resource Type", "Last Updated"})
		for _, entry := range entries {
			e, ok := entry.(map[string]interface{})
			if !ok {
				continue
			}
			res, ok := e["resource"].(map[string]interface{})
			if !ok {
				continue
			}

			id := getString(res, "id")
			rt := getString(res, "resourceType")
			meta, _ := res["meta"].(map[string]interface{})
			lastUpdated := getString(meta, "lastUpdated")

			table.Append([]string{id, rt, lastUpdated})
		}
	}

	table.Render()
	return nil
}

func printPatientBundleTable(entries []interface{}) error {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Name", "DOB", "Gender", "MRN"})

	for _, entry := range entries {
		e, ok := entry.(map[string]interface{})
		if !ok {
			continue
		}
		res, ok := e["resource"].(map[string]interface{})
		if !ok {
			continue
		}

		id := getString(res, "id")
		name := extractName(res)
		dob := getString(res, "birthDate")
		gender := getString(res, "gender")
		mrn := extractIdentifier(res, "MRN")

		table.Append([]string{id, name, dob, gender, mrn})
	}

	table.Render()
	return nil
}

func printPatientTable(patient map[string]interface{}) error {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Field", "Value"})
	table.SetAlignment(tablewriter.ALIGN_LEFT)

	table.Append([]string{"ID", getString(patient, "id")})
	table.Append([]string{"Name", extractName(patient)})
	table.Append([]string{"Birth Date", getString(patient, "birthDate")})
	table.Append([]string{"Gender", getString(patient, "gender")})
	table.Append([]string{"MRN", extractIdentifier(patient, "MRN")})

	if address := extractAddress(patient); address != "" {
		table.Append([]string{"Address", address})
	}

	if phone := extractTelecom(patient, "phone"); phone != "" {
		table.Append([]string{"Phone", phone})
	}

	if email := extractTelecom(patient, "email"); email != "" {
		table.Append([]string{"Email", email})
	}

	table.Render()
	return nil
}

func printObservationBundleTable(entries []interface{}) error {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Code", "Value", "Date", "Status"})

	for _, entry := range entries {
		e, ok := entry.(map[string]interface{})
		if !ok {
			continue
		}
		res, ok := e["resource"].(map[string]interface{})
		if !ok {
			continue
		}

		id := getString(res, "id")
		code := extractCodeText(res, "code")
		value := extractObservationValue(res)
		date := extractEffectiveDate(res)
		status := getString(res, "status")

		table.Append([]string{id, code, value, date, status})
	}

	table.Render()
	return nil
}

func printObservationTable(obs map[string]interface{}) error {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Field", "Value"})
	table.SetAlignment(tablewriter.ALIGN_LEFT)

	table.Append([]string{"ID", getString(obs, "id")})
	table.Append([]string{"Status", getString(obs, "status")})
	table.Append([]string{"Code", extractCodeText(obs, "code")})
	table.Append([]string{"Value", extractObservationValue(obs)})
	table.Append([]string{"Effective Date", extractEffectiveDate(obs)})
	table.Append([]string{"Category", extractCodeText(obs, "category")})

	table.Render()
	return nil
}

func printConditionBundleTable(entries []interface{}) error {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Condition", "Clinical Status", "Onset", "Recorded"})

	for _, entry := range entries {
		e, ok := entry.(map[string]interface{})
		if !ok {
			continue
		}
		res, ok := e["resource"].(map[string]interface{})
		if !ok {
			continue
		}

		id := getString(res, "id")
		code := extractCodeText(res, "code")
		status := extractCodeText(res, "clinicalStatus")
		onset := getString(res, "onsetDateTime")
		recorded := getString(res, "recordedDate")

		table.Append([]string{id, code, status, onset, recorded})
	}

	table.Render()
	return nil
}

func printConditionTable(cond map[string]interface{}) error {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Field", "Value"})
	table.SetAlignment(tablewriter.ALIGN_LEFT)

	table.Append([]string{"ID", getString(cond, "id")})
	table.Append([]string{"Condition", extractCodeText(cond, "code")})
	table.Append([]string{"Clinical Status", extractCodeText(cond, "clinicalStatus")})
	table.Append([]string{"Verification Status", extractCodeText(cond, "verificationStatus")})
	table.Append([]string{"Category", extractCodeText(cond, "category")})
	table.Append([]string{"Onset", getString(cond, "onsetDateTime")})
	table.Append([]string{"Recorded Date", getString(cond, "recordedDate")})

	table.Render()
	return nil
}

func printMedicationRequestBundleTable(entries []interface{}) error {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Medication", "Status", "Intent", "Authored"})

	for _, entry := range entries {
		e, ok := entry.(map[string]interface{})
		if !ok {
			continue
		}
		res, ok := e["resource"].(map[string]interface{})
		if !ok {
			continue
		}

		id := getString(res, "id")
		med := extractMedicationName(res)
		status := getString(res, "status")
		intent := getString(res, "intent")
		authored := getString(res, "authoredOn")

		table.Append([]string{id, med, status, intent, authored})
	}

	table.Render()
	return nil
}

func printMedicationRequestTable(med map[string]interface{}) error {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Field", "Value"})
	table.SetAlignment(tablewriter.ALIGN_LEFT)

	table.Append([]string{"ID", getString(med, "id")})
	table.Append([]string{"Medication", extractMedicationName(med)})
	table.Append([]string{"Status", getString(med, "status")})
	table.Append([]string{"Intent", getString(med, "intent")})
	table.Append([]string{"Authored On", getString(med, "authoredOn")})

	table.Render()
	return nil
}

func printEncounterBundleTable(entries []interface{}) error {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Type", "Status", "Class", "Period"})

	for _, entry := range entries {
		e, ok := entry.(map[string]interface{})
		if !ok {
			continue
		}
		res, ok := e["resource"].(map[string]interface{})
		if !ok {
			continue
		}

		id := getString(res, "id")
		encType := extractCodeText(res, "type")
		status := getString(res, "status")
		class := extractCoding(res, "class")
		period := extractPeriod(res)

		table.Append([]string{id, encType, status, class, period})
	}

	table.Render()
	return nil
}

func printEncounterTable(enc map[string]interface{}) error {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Field", "Value"})
	table.SetAlignment(tablewriter.ALIGN_LEFT)

	table.Append([]string{"ID", getString(enc, "id")})
	table.Append([]string{"Status", getString(enc, "status")})
	table.Append([]string{"Class", extractCoding(enc, "class")})
	table.Append([]string{"Type", extractCodeText(enc, "type")})
	table.Append([]string{"Period", extractPeriod(enc)})

	table.Render()
	return nil
}

func printGenericResourceTable(resource map[string]interface{}) error {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Field", "Value"})
	table.SetAlignment(tablewriter.ALIGN_LEFT)

	table.Append([]string{"Resource Type", getString(resource, "resourceType")})
	table.Append([]string{"ID", getString(resource, "id")})

	if meta, ok := resource["meta"].(map[string]interface{}); ok {
		table.Append([]string{"Last Updated", getString(meta, "lastUpdated")})
		table.Append([]string{"Version ID", getString(meta, "versionId")})
	}

	table.Render()
	return nil
}

// Helper functions

func getString(m map[string]interface{}, key string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	return ""
}

func extractName(patient map[string]interface{}) string {
	names, ok := patient["name"].([]interface{})
	if !ok || len(names) == 0 {
		return ""
	}

	name, ok := names[0].(map[string]interface{})
	if !ok {
		return ""
	}

	var parts []string

	if given, ok := name["given"].([]interface{}); ok {
		for _, g := range given {
			if s, ok := g.(string); ok {
				parts = append(parts, s)
			}
		}
	}

	if family, ok := name["family"].(string); ok {
		parts = append(parts, family)
	}

	if text, ok := name["text"].(string); ok && len(parts) == 0 {
		return text
	}

	return strings.Join(parts, " ")
}

func extractIdentifier(resource map[string]interface{}, typeCode string) string {
	identifiers, ok := resource["identifier"].([]interface{})
	if !ok {
		return ""
	}

	for _, id := range identifiers {
		identifier, ok := id.(map[string]interface{})
		if !ok {
			continue
		}

		if idType, ok := identifier["type"].(map[string]interface{}); ok {
			if text, ok := idType["text"].(string); ok && strings.Contains(strings.ToUpper(text), typeCode) {
				if value, ok := identifier["value"].(string); ok {
					return value
				}
			}
		}
	}

	// Return first identifier if specific type not found
	if len(identifiers) > 0 {
		if id, ok := identifiers[0].(map[string]interface{}); ok {
			if value, ok := id["value"].(string); ok {
				return value
			}
		}
	}

	return ""
}

func extractAddress(patient map[string]interface{}) string {
	addresses, ok := patient["address"].([]interface{})
	if !ok || len(addresses) == 0 {
		return ""
	}

	addr, ok := addresses[0].(map[string]interface{})
	if !ok {
		return ""
	}

	if text, ok := addr["text"].(string); ok {
		return text
	}

	var parts []string

	if lines, ok := addr["line"].([]interface{}); ok {
		for _, line := range lines {
			if s, ok := line.(string); ok {
				parts = append(parts, s)
			}
		}
	}

	if city, ok := addr["city"].(string); ok {
		parts = append(parts, city)
	}

	if state, ok := addr["state"].(string); ok {
		parts = append(parts, state)
	}

	if postal, ok := addr["postalCode"].(string); ok {
		parts = append(parts, postal)
	}

	return strings.Join(parts, ", ")
}

func extractTelecom(patient map[string]interface{}, system string) string {
	telecoms, ok := patient["telecom"].([]interface{})
	if !ok {
		return ""
	}

	for _, t := range telecoms {
		telecom, ok := t.(map[string]interface{})
		if !ok {
			continue
		}

		if sys, ok := telecom["system"].(string); ok && sys == system {
			if value, ok := telecom["value"].(string); ok {
				return value
			}
		}
	}

	return ""
}

func extractCodeText(resource map[string]interface{}, field string) string {
	code, ok := resource[field].(map[string]interface{})
	if !ok {
		// Try as array (for category)
		if codes, ok := resource[field].([]interface{}); ok && len(codes) > 0 {
			if c, ok := codes[0].(map[string]interface{}); ok {
				code = c
			}
		}
	}

	if code == nil {
		return ""
	}

	if text, ok := code["text"].(string); ok {
		return text
	}

	if codings, ok := code["coding"].([]interface{}); ok && len(codings) > 0 {
		if coding, ok := codings[0].(map[string]interface{}); ok {
			if display, ok := coding["display"].(string); ok {
				return display
			}
			if codeValue, ok := coding["code"].(string); ok {
				return codeValue
			}
		}
	}

	return ""
}

func extractCoding(resource map[string]interface{}, field string) string {
	coding, ok := resource[field].(map[string]interface{})
	if !ok {
		return ""
	}

	if display, ok := coding["display"].(string); ok {
		return display
	}

	if code, ok := coding["code"].(string); ok {
		return code
	}

	return ""
}

func extractObservationValue(obs map[string]interface{}) string {
	// valueQuantity
	if vq, ok := obs["valueQuantity"].(map[string]interface{}); ok {
		value := ""
		if v, ok := vq["value"].(float64); ok {
			value = fmt.Sprintf("%.2f", v)
		}
		if unit, ok := vq["unit"].(string); ok {
			value += " " + unit
		}
		return strings.TrimSpace(value)
	}

	// valueString
	if vs, ok := obs["valueString"].(string); ok {
		return vs
	}

	// valueCodeableConcept
	if vc, ok := obs["valueCodeableConcept"].(map[string]interface{}); ok {
		if text, ok := vc["text"].(string); ok {
			return text
		}
		if codings, ok := vc["coding"].([]interface{}); ok && len(codings) > 0 {
			if coding, ok := codings[0].(map[string]interface{}); ok {
				if display, ok := coding["display"].(string); ok {
					return display
				}
			}
		}
	}

	// valueBoolean
	if vb, ok := obs["valueBoolean"].(bool); ok {
		if vb {
			return "true"
		}
		return "false"
	}

	// valueInteger
	if vi, ok := obs["valueInteger"].(float64); ok {
		return fmt.Sprintf("%d", int(vi))
	}

	return ""
}

func extractEffectiveDate(obs map[string]interface{}) string {
	if ed, ok := obs["effectiveDateTime"].(string); ok {
		return ed
	}

	if ep, ok := obs["effectivePeriod"].(map[string]interface{}); ok {
		start := getString(ep, "start")
		end := getString(ep, "end")
		if start != "" && end != "" {
			return start + " - " + end
		}
		if start != "" {
			return start
		}
	}

	return ""
}

func extractMedicationName(med map[string]interface{}) string {
	// medicationCodeableConcept
	if mc, ok := med["medicationCodeableConcept"].(map[string]interface{}); ok {
		if text, ok := mc["text"].(string); ok {
			return text
		}
		if codings, ok := mc["coding"].([]interface{}); ok && len(codings) > 0 {
			if coding, ok := codings[0].(map[string]interface{}); ok {
				if display, ok := coding["display"].(string); ok {
					return display
				}
			}
		}
	}

	// medicationReference
	if mr, ok := med["medicationReference"].(map[string]interface{}); ok {
		if display, ok := mr["display"].(string); ok {
			return display
		}
		if ref, ok := mr["reference"].(string); ok {
			return ref
		}
	}

	return ""
}

func extractPeriod(resource map[string]interface{}) string {
	period, ok := resource["period"].(map[string]interface{})
	if !ok {
		return ""
	}

	start := getString(period, "start")
	end := getString(period, "end")

	if start != "" && end != "" {
		return start + " - " + end
	}
	if start != "" {
		return "Started: " + start
	}

	return ""
}
