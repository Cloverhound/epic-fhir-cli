package fhir

// Common FHIR resource types supported by Epic
const (
	ResourcePatient              = "Patient"
	ResourceObservation          = "Observation"
	ResourceCondition            = "Condition"
	ResourceMedicationRequest    = "MedicationRequest"
	ResourceMedicationStatement  = "MedicationStatement"
	ResourceAllergyIntolerance   = "AllergyIntolerance"
	ResourceProcedure            = "Procedure"
	ResourceDiagnosticReport     = "DiagnosticReport"
	ResourceEncounter            = "Encounter"
	ResourceImmunization         = "Immunization"
	ResourceCarePlan             = "CarePlan"
	ResourceCareTeam             = "CareTeam"
	ResourceGoal                 = "Goal"
	ResourceDocumentReference    = "DocumentReference"
	ResourcePractitioner         = "Practitioner"
	ResourcePractitionerRole     = "PractitionerRole"
	ResourceOrganization         = "Organization"
	ResourceLocation             = "Location"
	ResourceDevice               = "Device"
	ResourceAppointment          = "Appointment"
	ResourceSchedule             = "Schedule"
	ResourceSlot                 = "Slot"
	ResourceCoverage             = "Coverage"
	ResourceClaim                = "Claim"
	ResourceExplanationOfBenefit = "ExplanationOfBenefit"
	ResourceServiceRequest       = "ServiceRequest"
	ResourceSpecimen             = "Specimen"
	ResourceFamilyMemberHistory  = "FamilyMemberHistory"
	ResourceRelatedPerson        = "RelatedPerson"
	ResourceProvenance           = "Provenance"
	ResourceQuestionnaireResponse = "QuestionnaireResponse"
	ResourceConsent              = "Consent"
	ResourceMedicationAdministration = "MedicationAdministration"
	ResourceMedicationDispense   = "MedicationDispense"
	ResourceNutritionOrder       = "NutritionOrder"
	ResourceBinary               = "Binary"
)

// ResourceInfo contains metadata about a FHIR resource type
type ResourceInfo struct {
	Name           string
	Description    string
	SearchParams   []string
	CommonFilters  map[string]string
}

// SupportedResources returns information about all supported resources
var SupportedResources = map[string]ResourceInfo{
	ResourcePatient: {
		Name:        "Patient",
		Description: "Demographics and administrative information about a patient",
		SearchParams: []string{
			"_id", "identifier", "family", "given", "name", "birthdate",
			"gender", "address", "address-city", "address-state", "address-postalcode",
			"phone", "email", "telecom",
		},
	},
	ResourceObservation: {
		Name:        "Observation",
		Description: "Measurements and simple assertions about a patient",
		SearchParams: []string{
			"_id", "patient", "subject", "category", "code", "date",
			"status", "value-quantity", "component-code",
		},
		CommonFilters: map[string]string{
			"vital-signs":  "category=vital-signs",
			"laboratory":   "category=laboratory",
			"social-history": "category=social-history",
		},
	},
	ResourceCondition: {
		Name:        "Condition",
		Description: "Detailed information about conditions, problems, or diagnoses",
		SearchParams: []string{
			"_id", "patient", "subject", "category", "clinical-status",
			"code", "onset-date", "recorded-date", "verification-status",
		},
	},
	ResourceMedicationRequest: {
		Name:        "MedicationRequest",
		Description: "An order or request for medication",
		SearchParams: []string{
			"_id", "patient", "subject", "status", "intent",
			"medication", "authoredon", "requester",
		},
	},
	ResourceAllergyIntolerance: {
		Name:        "AllergyIntolerance",
		Description: "Allergy or intolerance information",
		SearchParams: []string{
			"_id", "patient", "clinical-status", "verification-status",
			"type", "category", "criticality", "code",
		},
	},
	ResourceProcedure: {
		Name:        "Procedure",
		Description: "Actions taken on or for a patient",
		SearchParams: []string{
			"_id", "patient", "subject", "date", "status",
			"code", "category", "performer",
		},
	},
	ResourceDiagnosticReport: {
		Name:        "DiagnosticReport",
		Description: "Diagnostic report findings and interpretation",
		SearchParams: []string{
			"_id", "patient", "subject", "category", "code",
			"date", "status", "result",
		},
	},
	ResourceEncounter: {
		Name:        "Encounter",
		Description: "An interaction between patient and healthcare provider",
		SearchParams: []string{
			"_id", "patient", "subject", "date", "status",
			"class", "type", "participant", "location",
		},
	},
	ResourceImmunization: {
		Name:        "Immunization",
		Description: "Immunization event information",
		SearchParams: []string{
			"_id", "patient", "date", "status",
			"vaccine-code", "location", "performer",
		},
	},
	ResourceCarePlan: {
		Name:        "CarePlan",
		Description: "Healthcare plan for a patient",
		SearchParams: []string{
			"_id", "patient", "subject", "date", "status",
			"category", "intent",
		},
	},
	ResourceCareTeam: {
		Name:        "CareTeam",
		Description: "Team of practitioners caring for a patient",
		SearchParams: []string{
			"_id", "patient", "subject", "status",
			"category", "participant",
		},
	},
	ResourceGoal: {
		Name:        "Goal",
		Description: "Describes desired objectives for patient care",
		SearchParams: []string{
			"_id", "patient", "subject", "lifecycle-status",
			"category", "target-date",
		},
	},
	ResourceDocumentReference: {
		Name:        "DocumentReference",
		Description: "Reference to a document",
		SearchParams: []string{
			"_id", "patient", "subject", "date", "status",
			"type", "category", "period", "author",
		},
	},
	ResourcePractitioner: {
		Name:        "Practitioner",
		Description: "Healthcare provider information",
		SearchParams: []string{
			"_id", "identifier", "name", "family", "given",
			"active", "address",
		},
	},
	ResourceOrganization: {
		Name:        "Organization",
		Description: "A formally or informally recognized group",
		SearchParams: []string{
			"_id", "identifier", "name", "type",
			"address", "active",
		},
	},
	ResourceLocation: {
		Name:        "Location",
		Description: "Details about a physical place",
		SearchParams: []string{
			"_id", "identifier", "name", "type",
			"address", "status",
		},
	},
	ResourceAppointment: {
		Name:        "Appointment",
		Description: "A booking of a healthcare event",
		SearchParams: []string{
			"_id", "patient", "date", "status",
			"practitioner", "location", "service-type",
		},
	},
	ResourceCoverage: {
		Name:        "Coverage",
		Description: "Insurance or payment information",
		SearchParams: []string{
			"_id", "patient", "beneficiary", "status",
			"type", "payor",
		},
	},
	ResourceServiceRequest: {
		Name:        "ServiceRequest",
		Description: "A request for a service to be performed",
		SearchParams: []string{
			"_id", "patient", "subject", "status",
			"intent", "code", "authored",
		},
	},
	ResourceFamilyMemberHistory: {
		Name:        "FamilyMemberHistory",
		Description: "Information about a patient's family member's health history",
		SearchParams: []string{
			"_id", "patient", "status", "relationship",
			"code", "date",
		},
	},
	ResourceConsent: {
		Name:        "Consent",
		Description: "A healthcare consumer's policy choices",
		SearchParams: []string{
			"_id", "patient", "status", "category",
			"scope", "date",
		},
	},
	ResourceMedicationAdministration: {
		Name:        "MedicationAdministration",
		Description: "Event of medication being administered to a patient",
		SearchParams: []string{
			"_id", "patient", "subject", "status",
			"medication", "effective-time", "performer",
		},
	},
	ResourceMedicationDispense: {
		Name:        "MedicationDispense",
		Description: "Dispensing a medication to a patient",
		SearchParams: []string{
			"_id", "patient", "subject", "status",
			"medication", "whenhandedover", "performer",
		},
	},
}

// GetResourceInfo returns information about a resource type
func GetResourceInfo(resourceType string) (ResourceInfo, bool) {
	info, ok := SupportedResources[resourceType]
	return info, ok
}

// ListResourceTypes returns all supported resource type names
func ListResourceTypes() []string {
	types := make([]string, 0, len(SupportedResources))
	for name := range SupportedResources {
		types = append(types, name)
	}
	return types
}
