# fhir-cli

A command-line interface for querying Epic's FHIR R4 APIs.

## Features

- **OAuth2 Backend Authentication**: JWT-based authentication using RS384 signing
- **All FHIR R4 Resources**: Support for Patient, Observation, Condition, MedicationRequest, and 25+ more resources
- **Multiple Output Formats**: JSON, YAML, and formatted tables
- **Profile Management**: Multiple configuration profiles for different endpoints
- **Token Caching**: Automatic token caching and refresh

## Installation

### From Source

```bash
git clone https://github.com/jbogarin/fhir-cli.git
cd fhir-cli
make build
```

### Install to GOPATH

```bash
make install
```

## Quick Start

### 1. Configure credentials

```bash
fhir-cli config init
```

This will prompt you for:
- Profile name (default: sandbox)
- Client ID (from Epic app registration)
- Private key source (file path or environment variable)
- FHIR Base URL
- Token URL
- Scopes

### 2. Verify authentication

```bash
fhir-cli auth token
```

### 3. Query resources

```bash
# Search for patients
fhir-cli patient search --name "Smith" -o table

# Get patient by ID
fhir-cli patient get abc123

# Get observations for a patient
fhir-cli observation search --patient abc123 --category vital-signs

# Get active medications
fhir-cli medication active abc123 -o table
```

## Usage

### Global Flags

| Flag | Description |
|------|-------------|
| `-o, --output` | Output format: json, table, yaml (default: json) |
| `-p, --profile` | Configuration profile to use (default: default) |
| `--config` | Custom config file path |

### Commands

#### Configuration

```bash
fhir-cli config init          # Interactive setup
fhir-cli config show          # Show current config
fhir-cli config set <key> <value>  # Set config value
fhir-cli config add-profile <name> # Add new profile
fhir-cli config use <profile>      # Set default profile
```

#### Authentication

```bash
fhir-cli auth token           # Get/refresh access token
fhir-cli auth status          # Check authentication status
fhir-cli auth logout          # Clear cached tokens
```

#### Resources

Each resource supports `get`, `search`, and often `list` commands:

```bash
# Patient
fhir-cli patient get <id>
fhir-cli patient search --name "Smith" --birthdate 1980-01-15
fhir-cli patient everything <id>

# Observation
fhir-cli observation get <id>
fhir-cli observation search --patient <id> --category vital-signs
fhir-cli observation vitals <patient-id>
fhir-cli observation labs <patient-id>

# Condition
fhir-cli condition get <id>
fhir-cli condition search --patient <id> --clinical-status active
fhir-cli condition active <patient-id>
fhir-cli condition problems <patient-id>

# Medication
fhir-cli medication get <id>
fhir-cli medication search --patient <id> --status active
fhir-cli medication active <patient-id>

# Allergy
fhir-cli allergy list <patient-id>
fhir-cli allergy search --patient <id> --criticality high

# Encounter
fhir-cli encounter list <patient-id>
fhir-cli encounter search --patient <id> --status finished

# And many more...
```

#### Generic Resource Access

```bash
# Get any resource by type and ID
fhir-cli resource get <ResourceType> <id>

# Search any resource
fhir-cli resource search <ResourceType> param=value

# Get server metadata
fhir-cli resource metadata

# List supported resources
fhir-cli resource list
```

## Supported Resources

| Resource | Commands |
|----------|----------|
| Patient | get, search, everything |
| Observation | get, search, vitals, labs |
| Condition | get, search, active, problems |
| MedicationRequest | get, search, active |
| AllergyIntolerance | get, search, list |
| Procedure | get, search, list |
| DiagnosticReport | get, search, labs, imaging |
| Encounter | get, search, list |
| Immunization | get, search, list |
| CarePlan | get, search, list |
| CareTeam | search |
| Goal | search |
| DocumentReference | get, search, list |
| Practitioner | get, search |
| PractitionerRole | search |
| Organization | get, search |
| Location | get, search |
| Appointment | get, search, list, upcoming |
| Schedule | get, search |
| Slot | search |
| Coverage | get, search, list |
| ExplanationOfBenefit | get, search |
| ServiceRequest | get, search |
| FamilyMemberHistory | get, search, list |
| RelatedPerson | get, search |
| Consent | get, search |
| Provenance | search |
| QuestionnaireResponse | get, search |

## Configuration

Configuration is stored in `~/.fhir-cli/config.yaml`:

```yaml
default: sandbox
profiles:
  sandbox:
    name: sandbox
    client_id: your-client-id
    private_key: env:FHIR_PRIVATE_KEY  # Read from environment variable
    token_url: https://fhir.epic.com/interconnect-fhir-oauth/oauth2/token
    fhir_base_url: https://fhir.epic.com/interconnect-fhir-oauth/api/FHIR/R4
    fhir_version: R4
    scopes: system/Patient.read system/Observation.read
    output_format: json

  production:
    name: production
    client_id: prod-client-id
    private_key: env:FHIR_PROD_PRIVATE_KEY
    token_url: https://epic.hospital.org/oauth2/token
    fhir_base_url: https://epic.hospital.org/api/FHIR/R4
    fhir_version: R4
    scopes: system/Patient.read system/Observation.read
```

### Private Key Configuration

The `private_key` field supports two formats:

| Format | Example | Description |
|--------|---------|-------------|
| Environment variable | `env:FHIR_PRIVATE_KEY` | Reads PEM-encoded key from the specified env var |
| File path | `/path/to/key.pem` | Reads key from file (supports `~` expansion) |

**Using environment variables (recommended for production):**

```bash
# Set the private key in your environment
export FHIR_PRIVATE_KEY="-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEA...
-----END RSA PRIVATE KEY-----"

# Or load from a file into the env var
export FHIR_PRIVATE_KEY=$(cat /secure/path/to/key.pem)

# Configure the CLI to use the env var
fhir-cli config set private_key "env:FHIR_PRIVATE_KEY"
```

### Multiple Environments

```bash
# Add profiles for different environments
fhir-cli config add-profile staging
fhir-cli config set client_id "staging-client" -p staging
fhir-cli config set private_key "env:FHIR_STAGING_KEY" -p staging
fhir-cli config set fhir_base_url "https://staging.hospital.org/api/FHIR/R4" -p staging

# Switch default environment
fhir-cli config use staging

# Or use -p flag for one-off commands
fhir-cli patient search --name Smith -p production
```

## Epic Sandbox

For testing, use Epic's sandbox endpoints:

| Version | Base URL |
|---------|----------|
| R4 | `https://fhir.epic.com/interconnect-fhir-oauth/api/FHIR/R4` |
| STU3 | `https://fhir.epic.com/interconnect-fhir-oauth/api/FHIR/STU3` |

## License

MIT
