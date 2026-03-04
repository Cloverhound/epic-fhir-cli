---
name: fhir-cli
description: "FHIR CLI: query and manage Epic FHIR R4 resources via the `fhir-cli` command-line tool. Use for searching patients, retrieving clinical data (observations, conditions, medications, allergies, encounters, procedures), checking authentication status, and exploring FHIR server capabilities."
argument-hint: "[command or resource-name]"
allowed-tools: Bash, Read, Grep, Glob
user-invocable: true
---

# FHIR CLI Skill

This skill uses the `fhir-cli` CLI tool to interact with Epic's FHIR R4 APIs using OAuth2 backend authentication with JWT assertions.

**Binary path** (update to match your installation):
```bash
fhir-cli
```

## Setup

Install via the one-liner:
```bash
curl -fsSL https://raw.githubusercontent.com/Cloverhound/epic-fhir-cli/main/install.sh | sh
```

On Windows (PowerShell):
```powershell
irm https://raw.githubusercontent.com/Cloverhound/epic-fhir-cli/main/install.ps1 | iex
```

Or install from source:
```bash
git clone https://github.com/Cloverhound/epic-fhir-cli.git
cd epic-fhir-cli
make install
```

## Authentication

The CLI uses OAuth2 backend authentication with RS384 JWT assertions. Tokens are stored securely in the OS keyring (macOS Keychain, Windows Credential Manager, Linux Secret Service/D-Bus) and reused until 5 minutes before expiry.

### Configuration

```bash
fhir-cli config init                    # Interactive setup (creates ~/.fhir-cli/config.yaml)
fhir-cli config show                    # Show current configuration and profiles
fhir-cli config add-profile prod        # Add a profile with defaults
fhir-cli config add-profile prod -i     # Add a profile with interactive wizard
fhir-cli config use prod                # Switch default profile
fhir-cli config set client_id "abc123"  # Set a config value for current profile
```

Config keys: `client_id`, `private_key`, `fhir_base_url`, `token_url`, `fhir_version`, `scopes`, `output_format`

Private key supports file paths (`/path/to/key.pem`) or environment variables (`env:FHIR_PRIVATE_KEY`).

### Auth Commands

```bash
fhir-cli auth token          # Get or refresh access token
fhir-cli auth token --force  # Force token refresh
fhir-cli auth status         # Check authentication status
fhir-cli auth logout         # Clear cached tokens
fhir-cli auth debug          # Debug authentication configuration
fhir-cli auth test-key       # Test if private key can be loaded
fhir-cli auth show-jwt       # Display the JWT that would be sent
```

## Global Flags

```bash
-o, --output <format>   # Output format: json (default), table, yaml
-p, --profile <name>    # Configuration profile to use
-v, --verbose           # Verbose output (show HTTP requests/responses)
    --config <path>     # Config file path override
```

## Command Reference

### Patient

```bash
fhir-cli patient get <id>
fhir-cli patient search --family "Smith" --given "John"
fhir-cli patient search --name "Smith" --birthdate "1970-01-01" --gender "male"
fhir-cli patient search --identifier "MRN|12345"
fhir-cli patient search --address-city "Madison" --address-state "WI"
fhir-cli patient everything <id>    # $everything operation
```

Search flags: `--name`, `--family`, `--given`, `--birthdate`, `--gender`, `--identifier`, `--address`, `--address-city`, `--address-state`, `--address-postalcode`, `--phone`, `--email`, `--count`

### Observation

```bash
fhir-cli observation get <id>
fhir-cli observation search --patient <id> --category "vital-signs"
fhir-cli observation search --patient <id> --code "8867-4" --date "ge2024-01-01"
fhir-cli observation vitals <patient-id>           # Vital signs shortcut
fhir-cli observation vitals <patient-id> --date "ge2024-01-01"
fhir-cli observation labs <patient-id>             # Lab results shortcut
```

Search flags: `--patient`, `--category`, `--code`, `--date`, `--status`, `--count`

### Condition

```bash
fhir-cli condition get <id>
fhir-cli condition search --patient <id> --clinical-status "active"
fhir-cli condition search --patient <id> --category "problem-list-item"
fhir-cli condition active <patient-id>     # Active conditions shortcut
fhir-cli condition problems <patient-id>   # Problem list shortcut
```

Search flags: `--patient`, `--clinical-status`, `--category`, `--code`, `--onset-date`, `--count`

### Medication

```bash
fhir-cli medication get <id>
fhir-cli medication search --patient <id> --status "active"
fhir-cli medication active <patient-id>                    # Active meds shortcut
fhir-cli medication statement search --patient <id>        # MedicationStatement
fhir-cli medication administration search --patient <id>   # MedicationAdministration
fhir-cli medication dispense search --patient <id>         # MedicationDispense
```

Search flags: `--patient`, `--status`, `--intent`, `--authoredon`, `--count`

### Encounter

```bash
fhir-cli encounter get <id>
fhir-cli encounter search --patient <id> --status "finished" --date "ge2024-01-01"
fhir-cli encounter search --patient <id> --class "inpatient"
fhir-cli encounter list <patient-id>
```

Search flags: `--patient`, `--status`, `--class`, `--date`, `--type`, `--count`

### Procedure

```bash
fhir-cli procedure get <id>
fhir-cli procedure search --patient <id> --status "completed"
fhir-cli procedure search --patient <id> --date "ge2024-01-01"
fhir-cli procedure list <patient-id>
```

Search flags: `--patient`, `--status`, `--code`, `--date`, `--category`, `--count`

### Allergy / AllergyIntolerance

Aliases: `allergy`, `allergies`

```bash
fhir-cli allergy get <id>
fhir-cli allergy search --patient <id> --clinical-status "active"
fhir-cli allergy search --patient <id> --criticality "high"
fhir-cli allergy list <patient-id>     # All allergies for patient
```

Search flags: `--patient`, `--clinical-status`, `--criticality`, `--type`, `--category`, `--count`

### Immunization

Aliases: `immunization`, `vaccine`, `imm`

```bash
fhir-cli immunization get <id>
fhir-cli immunization search --patient <id> --status "completed"
fhir-cli immunization search --patient <id> --date "ge2024-01-01"
fhir-cli immunization list <patient-id>
```

Search flags: `--patient`, `--status`, `--vaccine-code`, `--date`, `--count`

### Appointment

Aliases: `appointment`, `appt`

```bash
fhir-cli appointment get <id>
fhir-cli appointment search --patient <id> --status "booked"
fhir-cli appointment search --patient <id> --date "ge2024-06-01"
fhir-cli appointment list <patient-id>
fhir-cli appointment upcoming <patient-id>   # Upcoming appointments
```

Search flags: `--patient`, `--status`, `--date`, `--practitioner`, `--location`, `--service-type`, `--count`

### Schedule

```bash
fhir-cli schedule get <id>
fhir-cli schedule search --actor <practitioner-id>
```

Search flags: `--actor`, `--date`, `--count`

### Slot

```bash
fhir-cli slot search --schedule <schedule-id> --status "free"
fhir-cli slot search --schedule <schedule-id> --start "ge2024-06-01"
```

Search flags: `--schedule`, `--status`, `--start`, `--count`

### CarePlan

Aliases: `careplan`, `care-plan`

```bash
fhir-cli careplan get <id>
fhir-cli careplan search --patient <id> --status "active"
fhir-cli careplan list <patient-id>
```

Search flags: `--patient`, `--status`, `--category`, `--date`, `--count`

### CareTeam

Aliases: `careteam`, `care-team`

```bash
fhir-cli careteam search --patient <id>
```

Search flags: `--patient`, `--status`, `--count`

### Goal

```bash
fhir-cli goal search --patient <id> --lifecycle-status "active"
```

Search flags: `--patient`, `--lifecycle-status`, `--count`

### Coverage (Insurance)

Aliases: `coverage`, `insurance`

```bash
fhir-cli coverage get <id>
fhir-cli coverage search --patient <id> --status "active"
fhir-cli coverage list <patient-id>
```

Search flags: `--patient`, `--beneficiary`, `--status`, `--type`, `--payor`, `--count`

### ExplanationOfBenefit

Aliases: `eob`, `explanation-of-benefit`

```bash
fhir-cli eob get <id>
fhir-cli eob search --patient <id>
```

Search flags: `--patient`, `--status`, `--created`, `--count`

### ServiceRequest

Aliases: `servicerequest`, `service-request`, `order`

```bash
fhir-cli servicerequest get <id>
fhir-cli servicerequest search --patient <id> --status "active"
```

Search flags: `--patient`, `--status`, `--intent`, `--category`, `--count`

### DiagnosticReport

Aliases: `diagnostic`, `dx`

```bash
fhir-cli diagnostic get <id>
fhir-cli diagnostic search --patient <id> --category "LAB"
fhir-cli diagnostic labs <patient-id>       # Lab reports shortcut
fhir-cli diagnostic imaging <patient-id>    # Imaging reports shortcut
```

Search flags: `--patient`, `--category`, `--code`, `--status`, `--date`, `--count`

### DocumentReference

Aliases: `document`, `doc`

```bash
fhir-cli document get <id>
fhir-cli document search --patient <id> --type "clinical-note"
fhir-cli document list <patient-id>
```

Search flags: `--patient`, `--status`, `--type`, `--category`, `--date`, `--period`, `--count`

### Binary

```bash
fhir-cli binary get <id>    # Get binary content by ID
```

### FamilyMemberHistory

Aliases: `familyhistory`, `family-history`, `fmh`

```bash
fhir-cli familyhistory get <id>
fhir-cli familyhistory search --patient <id>
fhir-cli familyhistory list <patient-id>
```

Search flags: `--patient`, `--status`, `--relationship`, `--code`, `--count`

### RelatedPerson

Aliases: `relatedperson`, `related-person`, `contact`

```bash
fhir-cli relatedperson get <id>
fhir-cli relatedperson search --patient <id>
```

Search flags: `--patient`, `--name`, `--relationship`, `--count`

### Consent

```bash
fhir-cli consent get <id>
fhir-cli consent search --patient <id> --status "active"
```

Search flags: `--patient`, `--status`, `--category`, `--count`

### Provenance

```bash
fhir-cli provenance search --target <resource-reference>
fhir-cli provenance search --patient <id>
```

Search flags: `--target`, `--patient`, `--recorded`, `--count`

### QuestionnaireResponse

Aliases: `questionnaireresponse`, `questionnaire-response`, `qr`

```bash
fhir-cli questionnaireresponse get <id>
fhir-cli questionnaireresponse search --patient <id>
```

Search flags: `--patient`, `--status`, `--authored`, `--count`

### Practitioner

Aliases: `practitioner`, `provider`

```bash
fhir-cli practitioner get <id>
fhir-cli practitioner search --family "Smith" --given "Jane"
fhir-cli practitioner search --name "Smith"
fhir-cli practitioner role search --specialty "cardiology"
```

Search flags: `--name`, `--family`, `--given`, `--identifier`, `--active`, `--count`
Role search flags: `--practitioner`, `--specialty`, `--organization`, `--count`

### Organization

Aliases: `organization`, `org`

```bash
fhir-cli organization get <id>
fhir-cli organization search --name "General Hospital"
```

Search flags: `--name`, `--identifier`, `--type`, `--address`, `--count`

### Location

```bash
fhir-cli location get <id>
fhir-cli location search --name "Emergency"
```

Search flags: `--name`, `--identifier`, `--type`, `--address`, `--count`

### Generic Resource Operations

For any FHIR resource type not covered by specific commands:

```bash
fhir-cli resource get <ResourceType> <id>
fhir-cli resource search <ResourceType> param1=value1 param2=value2
fhir-cli resource list                    # List all supported resource types
fhir-cli resource metadata                # Get server capability statement
fhir-cli resource operation <operation>   # Execute a FHIR operation
fhir-cli resource operation <operation> --resource-type <type> --id <id>
```

### Version

```bash
fhir-cli version    # Print version number
```

## Epic-Specific Notes

- **Patient search**: Epic requires at least one search parameter. Use `--family` and `--given` instead of `--name` for more reliable results with Epic systems.
- **Date formats**: Use FHIR date prefixes: `ge2024-01-01` (>=), `le2024-12-31` (<=), `eq2024-06-15` (exact).
- **Identifiers**: Use the format `system|value`, e.g., `--identifier "MRN|12345"`.
- **Scopes**: The OAuth2 scopes in your profile determine which resources you can access. Add scopes like `system/Immunization.read` to access additional resources.
- **Sandbox**: Epic's public sandbox at `https://fhir.epic.com/interconnect-fhir-oauth/api/FHIR/R4` has limited test data. Register at https://fhir.epic.com for credentials.
- **Token caching**: Tokens are securely stored in the OS keyring and automatically reused. Use `fhir-cli auth token --force` to force a refresh, or `fhir-cli auth logout` to clear stored tokens from the keyring.

## Output Handling

```bash
# Pretty JSON (default)
fhir-cli patient search --family "Smith"

# Table format
fhir-cli patient search --family "Smith" -o table

# YAML format
fhir-cli patient search --family "Smith" -o yaml

# Save to file for processing
fhir-cli patient search --family "Smith" > /tmp/patients.json
```

## When Answering Questions

1. **Check config first** with `fhir-cli config show` to confirm a profile is configured
2. **Check auth** with `fhir-cli auth status` to confirm valid authentication
3. **Use verbose mode** (`-v`) for debugging failed requests to see HTTP details
4. **Prefer `--family`/`--given`** over `--name` for patient searches on Epic systems
5. **Use shortcut commands** like `observation vitals`, `condition active`, `medication active` when appropriate
6. **Write to temp files** when gathering data for analysis, then read the files
7. **Use `--count`** to control result set size (defaults vary by resource)
